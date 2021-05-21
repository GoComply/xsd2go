package xsd

import "encoding/xml"

// Attribute defines single XML attribute
type Documentation struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2001/XMLSchema documentation"`
	Source   string   `xml:"source,attr"`
	InnerXml string   `xml:",innerxml"`
}
