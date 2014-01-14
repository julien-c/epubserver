package epubserver

import (
	"encoding/xml"
)

// XML Parsing
//
type Rootfile struct {
	Path      string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Container struct {
	XMLName   xml.Name   `xml:"container"`
	Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
}
