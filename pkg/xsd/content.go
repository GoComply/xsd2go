package xsd

import (
	"encoding/xml"
)

type GenericContent interface {
	Attributes() []Attribute
	compile(*Schema)
}
type SimpleContent struct {
	XMLName   xml.Name   `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
	Extension *Extension `xml:"extension"`
}

func (sc *SimpleContent) Attributes() []Attribute {
	if sc.Extension != nil {
		return sc.Extension.Attributes
	}
	return []Attribute{}
}

func (sc *SimpleContent) compile(sch *Schema) {
}

type Extension struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base       string      `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
}

type ComplexContent struct {
	XMLName   xml.Name   `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension *Extension `xml:"extension"`
	schema    *Schema    `xml:"-"`
}

func (cc *ComplexContent) Attributes() []Attribute {
	if cc.Extension != nil {
		return cc.Extension.Attributes
	}
	return []Attribute{}
}

func (c *ComplexContent) compile(sch *Schema) {
	c.schema = sch
}
