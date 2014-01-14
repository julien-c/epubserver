package epubserver

import (
	"encoding/xml"
)

// XML Parsing
//
type Package struct {
	XMLName  xml.Name `xml:"package"`
	Version  string   `xml:"version,attr"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Manifest struct {
	Items []ManifestItem `xml:"item"`
}

type ManifestItem struct {
	Href      string `xml:"href,attr"`
	Id        string `xml:"id,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	Toc      string    `xml:"toc,attr"`
	Itemrefs []Itemref `xml:"itemref"`
}

type Itemref struct {
	Idref  string `xml:"idref,attr"`
	Linear string `xml:"linear,attr"`
}
