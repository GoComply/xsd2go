package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

type Type interface {
	GoName() string
	GoTypeName() string
	Schema() *Schema
	Attributes() []Attribute
	Elements() []Element
}

type ComplexType struct {
	XMLName          xml.Name       `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name             string         `xml:"name,attr"`
	Mixed            string         `xml:"mixed,attr"`
	AttributesDirect []Attribute    `xml:"attribute"`
	Sequence         *Sequence      `xml:"sequence"`
	schema           *Schema        `xml:"-"`
	SimpleContent    *SimpleContent `xml:"simpleContent"`
}

func (ct *ComplexType) Attributes() []Attribute {
	return ct.AttributesDirect
}

func (ct *ComplexType) Elements() []Element {
	if ct.Sequence != nil {
		return ct.Sequence.Elements
	}
	return []Element{}
}

func (ct *ComplexType) GoName() string {
	return strcase.ToCamel(ct.Name)
}

func (ct *ComplexType) GoTypeName() string {
	return ct.GoName()
}

func (ct *ComplexType) Schema() *Schema {
	return ct.schema
}

func (ct *ComplexType) compile(sch *Schema) {
	ct.schema = sch
	if ct.Sequence != nil {
		ct.Sequence.compile(sch)
	}
}

type Sequence struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Elements []Element `xml:"element"`
}

func (s *Sequence) compile(sch *Schema) {
	for idx, _ := range s.Elements {
		el := &s.Elements[idx]
		el.compile(sch)
	}
}

type SimpleContent struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
}

type SimpleType struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Name    string   `xml:"name,attr"`
	schema  *Schema  `xml:"-"`
}

func (st *SimpleType) GoName() string {
	return strcase.ToCamel(st.Name)
}

func (st *SimpleType) GoTypeName() string {
	return "string"
}

func (st *SimpleType) Schema() *Schema {
	return st.schema
}

func (st *SimpleType) compile(sch *Schema) {
	st.schema = sch
}

func (st *SimpleType) Attributes() []Attribute {
	return []Attribute{}
}

func (st *SimpleType) Elements() []Element {
	return []Element{}
}

type staticType string

func (st staticType) GoName() string {
	return string(st)
}

func (ct staticType) GoTypeName() string {
	return ct.GoName()
}

func (st staticType) Attributes() []Attribute {
	return []Attribute{}
}

func (st staticType) Elements() []Element {
	return []Element{}
}

func (st staticType) Schema() *Schema {
	return nil
}

func StaticType(name string) staticType {
	if name == "string" || name == "dateTime" || name == "base64Binary" || name == "normalizedString" {
		return staticType("string")
	} else if name == "decimal" {
		return "float64"
	}
	panic("Type xsd:" + name + " not implemented")
	return staticType(name)
}
