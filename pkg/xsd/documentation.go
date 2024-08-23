package xsd

import "encoding/xml"

type Documentation struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema documentation"`
	Source  string   `xml:"source,attr"`
	Lang    string   `xml:"xml:lang,attr"`
	Content []byte   `xml:",innerxml"`
}
