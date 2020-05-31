package xsd

import (
	"encoding/xml"
)

type GenericContent interface {
	Attributes() []Attribute
	Elements() []Element
	compile(*Schema, *Element)
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

func (sc *SimpleContent) Elements() []Element {
	if sc.Extension != nil {
		return sc.Extension.Elements()
	}
	return []Element{}
}

func (sc *SimpleContent) compile(sch *Schema, parentElement *Element) {
}

type Extension struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base       string      `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
	Sequence   *Sequence   `xml:"sequence"`
}

func (ext *Extension) Elements() []Element {
	if ext.Sequence != nil {
		l := ext.Sequence.Elements()
		return l
	}
	return []Element{}
}

func (ext *Extension) compile(sch *Schema, parentElement *Element) {
	if ext.Sequence != nil {
		ext.Sequence.compile(sch, parentElement)
	}

}

type ComplexContent struct {
	XMLName   xml.Name   `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension *Extension `xml:"extension"`
}

func (cc *ComplexContent) Attributes() []Attribute {
	if cc.Extension != nil {
		return cc.Extension.Attributes
	}
	return []Attribute{}
}

func (cc *ComplexContent) Elements() []Element {
	if cc.Extension != nil {
		return cc.Extension.Elements()
	}
	return []Element{}
}

func (c *ComplexContent) compile(sch *Schema, parentElement *Element) {
	if c.Extension != nil {
		c.Extension.compile(sch, parentElement)
	}
}
