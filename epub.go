package epubserver

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
)

type Epub struct {
	Zip      *zip.ReadCloser
	Manifest []ManifestItem
	Spine    []Itemref
	Opfdir   string
}



func (epub *Epub) readFile(name string) (io.ReadCloser, error) {
	for _, f := range epub.Zip.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			} else {
				return rc, nil
			}
		}
	}
	return nil, errors.New("File not found")
}

func (epub *Epub) getFile(name string) ([]byte, error) {
	rc, err := epub.readFile(name)
	
	if err != nil {
		return nil, err
	} else {
		defer rc.Close()
		return ioutil.ReadAll(rc)
	}
}

func (epub *Epub) getFileAsString(name string) (string, error) {
	data, err := epub.getFile(name)
	return string(data), err
}

func (epub *Epub) getItemFromId(id string) (ManifestItem, error) {
	for _, item := range epub.Manifest {
		if item.Id == id {
			return item, nil
		}
	}
	return ManifestItem{}, errors.New("Item not found")
}

func (epub *Epub) getItemFromHref(href string) (ManifestItem, error) {
	for _, item := range epub.Manifest {
		if epub.FullHref(item) == href {
			return item, nil
		}
	}
	return ManifestItem{}, errors.New("Item not found")
}

func (epub *Epub) getSpinePrevNext(id string) (isInSpine bool, prev, next string) {
	for pos, item := range epub.Spine {
		if item.Idref == id {
			var prev, next string
			if pos > 0 {
				prev = epub.Spine[pos-1].Idref
			}
			if pos < len(epub.Spine) - 1 {
				next = epub.Spine[pos+1].Idref
			}
			return true, prev, next
		}
	}
	return false, "", ""
}


func (epub *Epub) FullHref(item ManifestItem) (string) {
	if epub.Opfdir == "." {
		return item.Href
	} else {
		return epub.Opfdir + "/" + item.Href
	}
}
