package xsd

import (
	"encoding/xml"
)

// Attribute defines single XML attribute
type Attribute struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Use     string   `xml:"use,attr"`
}
