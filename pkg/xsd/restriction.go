package xsd

import (
	"encoding/xml"
)

type Restriction struct {
	XMLName          xml.Name       `xml:"http://www.w3.org/2001/XMLSchema restriction"`
	Base             reference      `xml:"base,attr"`
	AttributesDirect []Attribute    `xml:"attribute"`
	SimpleContent    *SimpleContent `xml:"simpleContent"`
	typ              Type
}

func (r *Restriction) compile(sch *Schema, parentElement *Element) {
	for idx, _ := range r.AttributesDirect {
		attribute := &r.AttributesDirect[idx]
		attribute.compile(sch)
	}
	if r.SimpleContent != nil {
		r.SimpleContent.compile(sch, parentElement)
	}

	if r.Base == "" {
		panic("Not implemented: xsd:extension/@base empty, cannot extend unknown type")
	}

	r.typ = sch.findReferencedType(r.Base)
	if r.typ == nil {
		panic("Cannot build xsd:extension: unknown type: " + string(r.Base))
	}
	r.typ.compile(sch, parentElement)
}

func (r *Restriction) Attributes() []Attribute {
	result := make([]Attribute, 0)
	if r.typ != nil {
		result = append(result, r.typ.Attributes()...)
	}
	if r.SimpleContent != nil {
		result = append(result, r.SimpleContent.Attributes()...)
	}
	result = deduplicateAttributes(append(result, r.AttributesDirect...))

	return result
}
