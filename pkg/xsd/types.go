package xsd

import (
	"encoding/xml"
)

type ComplexType struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name       string      `xml:"name,attr"`
	Mixed      string      `xml:"mixed,attr"`
	Attributes []Attribute `xml:"attribute"`
}
