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
	if ct.SimpleContent != nil {
		return ct.SimpleContent.Attributes()
	}
	return ct.AttributesDirect
}

func (ct *ComplexType) Elements() []Element {
	if ct.Sequence != nil {
		return ct.Sequence.Elements()
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

	// Handle improbable name clash. Consider XSD defining two attributes on the element:
	// "id" and "Id", this would create name clash given the camelization we do.
	goNames := map[string]uint{}
	for idx, _ := range ct.Attributes() {
		attribute := &ct.Attributes()[idx]
		attribute.compile(sch)

		count := goNames[attribute.GoName()]
		count += 1
		goNames[attribute.GoName()] = count
		attribute.DuplicateCount = count
		// Second GoName may be different depending on the DuplicateCount
		goNames[attribute.GoName()] = count
	}

	if ct.SimpleContent != nil && len(ct.AttributesDirect) > 1 {
		panic("Not implemented: complexType " + ct.Name + " defines direct attribute and simple content")
	}
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

type Extension struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base       string      `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
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
