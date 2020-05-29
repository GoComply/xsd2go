package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

type Type interface {
	GoName() string
	Schema() *Schema
}

type ComplexType struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name       string      `xml:"name,attr"`
	Mixed      string      `xml:"mixed,attr"`
	Attributes []Attribute `xml:"attribute"`
	Sequence   *Sequence   `xml:"sequence"`
	schema     *Schema     `xml:"-"`
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

func (ct *ComplexType) Schema() *Schema {
	return ct.schema
}

func (ct *ComplexType) compile(sch *Schema) {
	ct.schema = sch
}

type Sequence struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Elements []Element `xml:"element"`
}

type SimpleType struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Name    string   `xml:"name,attr"`
	schema  *Schema  `xml:"-"`
}

func (st *SimpleType) GoName() string {
	return strcase.ToCamel(st.Name)
}

func (st *SimpleType) Schema() *Schema {
	return st.schema
}

func (st *SimpleType) compile(sch *Schema) {
	st.schema = sch
}

type staticType string

func (st staticType) GoName() string {
	return string(st)
}

func (st staticType) Schema() *Schema {
	return nil
}

func StaticType(name string) staticType {
	if name == "string" || name == "dateTime" || name == "base64Binary" || name == "normalizedString" {
		return staticType("string")
	} else if name == "decimal" {
		return "uint64"
	}
	panic("Type xsd:" + name + " not implemented")
	return staticType(name)
}
