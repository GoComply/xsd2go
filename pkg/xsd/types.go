package xsd

import (
	"encoding/xml"
)

type ComplexType struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name       string      `xml:"name,attr"`
	Mixed      string      `xml:"mixed,attr"`
	Attributes []Attribute `xml:"attribute"`
	Sequence   *Sequence   `xml:"sequence"`
}

func (ct *ComplexType) Elements() []Element {
	if ct.Sequence != nil {
		return ct.Sequence.Elements
	}
	return []Element{}
}

type Sequence struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Elements []Element `xml:"element"`
}
