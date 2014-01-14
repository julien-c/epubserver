package epubserver

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

func Open(filename string) (*Epub, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	// We need r to stay open forever

	epub := &Epub{Zip: r}

	mimetype, _ := epub.getFileAsString("mimetype")
	if mimetype != "application/epub+zip" {
		return new(Epub), errors.New("Unsupported mime type.")
	}

	containerFile, _ := epub.getFile("META-INF/container.xml")
	c := new(Container)

	err = xml.Unmarshal(containerFile, c)
	if err != nil {
		return nil, err
	}

	rootFile, _ := epub.getFile(c.Rootfiles[0].Path)
	epub.Opfdir = path.Dir(c.Rootfiles[0].Path)

	pkg := new(Package)
	err = xml.Unmarshal(rootFile, pkg)
	if err != nil {
		return nil, err
	}
	epub.Manifest = pkg.Manifest.Items
	epub.Spine = pkg.Spine.Itemrefs

	return epub, nil
}

func (epub *Epub) Serve() {
	http.Handle("/javascripts/", http.FileServer(http.Dir("../public")))
	http.HandleFunc("/", epub.ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func (epub *Epub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path[1:]
	if reqPath == "" {
		// Redirect to first component:
		firstComponent := epub.Spine[0].Idref
		item, _ := epub.getItemFromId(firstComponent)
		http.Redirect(w, r, epub.FullHref(item), 302)
	} else {
		item, err := epub.getItemFromHref(reqPath)
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
			return
		}
		file, err := epub.getFile(reqPath)

		isInSpine, prev, next := epub.getSpinePrevNext(item.Id)
		if isInSpine {
			var prevHref, nextHref string
			if prev != "" {
				prevItem, _ := epub.getItemFromId(prev)
				prevHref = epub.FullHref(prevItem)
			}
			if next != "" {
				nextItem, _ := epub.getItemFromId(next)
				nextHref = epub.FullHref(nextItem)
			}
			js := fmt.Sprintf("<script>var componentPrev = %q; var componentNext = %q; %s</script>", prevHref, nextHref, ReaderJs)
			fileString := strings.Replace(string(file), "</body>", js+"</body>", 1)
			w.Header().Set("Content-Type", item.MediaType)
			w.Write([]byte(fileString))
		} else {
			w.Header().Set("Content-Type", item.MediaType)
			w.Write(file)
		}
	}
}
