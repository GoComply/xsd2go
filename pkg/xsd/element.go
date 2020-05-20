package xsd

import (
	"encoding/xml"
)

// Element defines single XML element
type Element struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name        string       `xml:"name,attr"`
	ComplexType *ComplexType `xml:"complexType"`
}

func (e *Element) Attributes() []Attribute {
	if e.ComplexType != nil {
		return e.ComplexType.Attributes
	}
	return []Attribute{}
}
