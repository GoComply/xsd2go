package xsd

import (
	"encoding/xml"
)

type Restriction struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema restriction"`
	Base       string      `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
}

func (r *Restriction) compile(sch *Schema) {
	for idx, _ := range r.Attributes {
		attribute := &r.Attributes[idx]
		attribute.compile(sch)
	}
}
