package xsd

import (
	"encoding/xml"
)

type GenericContent interface {
	Attributes() []Attribute
	Elements() []Element
	ContainsText() bool
	compile(*Schema, *Element)
}
type SimpleContent struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
	Extension   *Extension   `xml:"extension"`
	Restriction *Restriction `xml:"restriction"`
}

func (sc *SimpleContent) Attributes() []Attribute {
	if sc.Extension != nil {
		return sc.Extension.Attributes()
	} else if sc.Restriction != nil {
		return sc.Restriction.Attributes()
	}
	return []Attribute{}
}

func (sc *SimpleContent) ContainsText() bool {
	return sc.Extension != nil && sc.Extension.ContainsText()
}

func (sc *SimpleContent) Elements() []Element {
	if sc.Extension != nil {
		return sc.Extension.Elements()
	}
	return []Element{}
}

func (sc *SimpleContent) compile(sch *Schema, parentElement *Element) {
	if sc.Extension != nil {
		sc.Extension.compile(sch, parentElement)
	}
	if sc.Restriction != nil {
		sc.Restriction.compile(sch, parentElement)
	}
}

type ComplexContent struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension   *Extension   `xml:"extension"`
	Restriction *Restriction `xml:"restriction"`
}

func (cc *ComplexContent) Attributes() []Attribute {
	if cc.Extension != nil {
		return cc.Extension.Attributes()
	} else if cc.Restriction != nil {
		return cc.Restriction.Attributes()
	}
	return []Attribute{}
}

func (cc *ComplexContent) Elements() []Element {
	if cc.Extension != nil {
		return cc.Extension.Elements()
	}
	return []Element{}
}

func (cc *ComplexContent) ContainsText() bool {
	return cc.Extension != nil && cc.Extension.ContainsText()
}

func (c *ComplexContent) compile(sch *Schema, parentElement *Element) {
	if c.Extension != nil {
		c.Extension.compile(sch, parentElement)
	}
	if c.Restriction != nil {
		if c.Extension != nil {
			panic("Not implemented: xsd:complexContent defines xsd:restriction and xsd:extension")
		}
		c.Restriction.compile(sch, parentElement)
	}
}
