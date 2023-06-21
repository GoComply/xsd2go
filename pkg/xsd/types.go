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
	ContainsText() bool
	IncludeTypeTemplate() bool
	IncludeElementTemplate() bool
	IncludeComplexTypeTemplate() bool
	IncludeTemplateName() string
	compile(*Schema, *Element)
}

func injectSchemaIntoAttributes(schema *Schema, intermAttributes []Attribute) []Attribute {
	attributesWithProperScema := make([]Attribute, len(intermAttributes))
	for idx, attribute := range intermAttributes {
		attribute.schema = schema
		attributesWithProperScema[idx] = attribute
	}
	return attributesWithProperScema
}

type ComplexType struct {
	XMLName          xml.Name        `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name             string          `xml:"name,attr"`
	Mixed            bool            `xml:"mixed,attr"`
	AttributesDirect []Attribute     `xml:"attribute"`
	Sequence         *Sequence       `xml:"sequence"`
	SequenceAll      *SequenceAll    `xml:"all"`
	schema           *Schema         `xml:"-"`
	SimpleContent    *SimpleContent  `xml:"simpleContent"`
	ComplexContent   *ComplexContent `xml:"complexContent"`
	Choice           *Choice         `xml:"choice"`
	content          GenericContent  `xml:"-"`
	override         Override        `xml:"-"`
}

func (ct *ComplexType) Attributes() []Attribute {
	if ct.content != nil {
		return ct.content.Attributes()
	}
	return ct.AttributesDirect
}

func (ct *ComplexType) Elements() []Element {
	if ct.Sequence != nil {
		return ct.Sequence.Elements()
	} else if ct.SequenceAll != nil {
		return ct.SequenceAll.Elements()
	} else if ct.content != nil {
		return ct.content.Elements()
	} else if ct.Choice != nil {
		return ct.Choice.Elements()
	}
	return []Element{}
}

func (ct *ComplexType) GoName() string {
	return strcase.ToCamel(ct.Name)
}

func (ct *ComplexType) GoTypeName() string {
	if ct.SimpleContent != nil && ct.SimpleContent.Extension != nil && ct.SimpleContent.Extension.ContainsText() {
		return ct.SimpleContent.Extension.typ.GoName()
	}
	return ct.GoName()
}

func (ct *ComplexType) ContainsText() bool {
	return ct.content != nil && ct.content.ContainsText()
}

func (ct *ComplexType) Schema() *Schema {
	return ct.schema
}

func (ct *ComplexType) IncludeTypeTemplate() bool {
	return ct.override.TemplateUsed && ct.override.IsIncl
}

func (ct *ComplexType) IncludeElementTemplate() bool {
	return ct.override.TemplateUsed && ct.override.IsElem
}

func (ct *ComplexType) IncludeComplexTypeTemplate() bool {
	return ct.override.TemplateUsed && ct.override.IsCompTyp
}

func (ct *ComplexType) IncludeTemplateName() string {
	return ct.override.TemplateName
}

func (ct *ComplexType) compile(sch *Schema, parentElement *Element) {
	ct.schema = sch
	if ct.Sequence != nil {
		ct.Sequence.compile(sch, parentElement)
	}
	if ct.SequenceAll != nil {
		if ct.Sequence != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:sequence and xsd:all")
		}
		ct.SequenceAll.compile(sch, parentElement)
	}

	// Handle improbable name clash. Consider XSD defining two attributes on the element:
	// "id" and "Id", this would create name clash given the camelization we do.
	goNames := map[string]uint{}
	for idx := range ct.Attributes() {
		attribute := &ct.Attributes()[idx]
		attribute.compile(sch)

		count := goNames[attribute.GoName()]
		count += 1
		goNames[attribute.GoName()] = count
		attribute.DuplicateCount = count
		// Second GoName may be different depending on the DuplicateCount
		goNames[attribute.GoName()] = count
	}

	if ct.ComplexContent != nil {
		ct.content = ct.ComplexContent
		if ct.SimpleContent != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:simpleContent and xsd:complexContent together")
		}
	} else if ct.SimpleContent != nil {
		ct.content = ct.SimpleContent
	}

	if ct.content != nil {
		if len(ct.AttributesDirect) > 1 {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines direct attribute and xsd:content")
		}
		if ct.Sequence != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:sequence and xsd:content")
		}
		if ct.SequenceAll != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:all and xsd:content")
		}
		ct.content.compile(sch, parentElement)
	}

	if ct.Choice != nil {
		if ct.content != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:choice and xsd:content")
		}
		if ct.Sequence != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:choice and xsd:sequence")
		}
		if ct.SequenceAll != nil {
			panic("Not implemented: xsd:complexType " + ct.Name + " defines xsd:all and xsd:sequence")
		}
		ct.Choice.compile(sch, parentElement)
	}

	if tmpl, found := sch.TemplateOverrides[ct.GoName()]; found {
		tmpl.TemplateUsed = true
		ct.override = tmpl
		sch.TemplateOverrides[ct.GoName()] = tmpl
	}
}

type SimpleType struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Name        string       `xml:"name,attr"`
	Restriction *Restriction `xml:"restriction"`
	schema      *Schema      `xml:"-"`
	override    Override     `xml:"-"`
}

func (st *SimpleType) GoName() string {
	return strcase.ToCamel(st.Name)
}

func (st *SimpleType) GoTypeName() string {
	if st.Restriction != nil && st.Restriction.typ != nil {
		return st.Restriction.typ.GoTypeName()
	}
	return "string"
}

func (st *SimpleType) Schema() *Schema {
	return st.schema
}

func (st *SimpleType) IncludeTypeTemplate() bool {
	return st.override.TemplateUsed && st.override.IsIncl
}

func (st *SimpleType) IncludeElementTemplate() bool {
	return st.override.TemplateUsed && st.override.IsElem
}

func (st *SimpleType) IncludeComplexTypeTemplate() bool {
	return false
}

func (st *SimpleType) IncludeTemplateName() string {
	return st.override.TemplateName
}

func (st *SimpleType) compile(sch *Schema, parentElement *Element) {
	if st.schema == nil {
		st.schema = sch
	}

	if st.Restriction != nil {
		st.Restriction.compile(sch, parentElement)
	}

	if tmpl, found := sch.TemplateOverrides[st.GoName()]; found {
		tmpl.TemplateUsed = true
		st.override = tmpl
		sch.TemplateOverrides[st.GoName()] = tmpl
	}
}

func (st *SimpleType) Attributes() []Attribute {
	return []Attribute{}
}

func (st *SimpleType) Elements() []Element {
	return []Element{}
}

func (st *SimpleType) Enums() []Enumeration {
	if st.Restriction != nil {
		return st.Restriction.Enums()
	}
	return []Enumeration{}
}

func (st *SimpleType) ContainsText() bool {
	return true
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

func (staticType) ContainsText() bool {
	return true
}

func (st staticType) IncludeTypeTemplate() bool {
	return false
}

func (st staticType) IncludeElementTemplate() bool {
	return false
}

func (st staticType) IncludeComplexTypeTemplate() bool {
	return false
}

func (st staticType) IncludeTemplateName() string {
	return ""
}

func (st staticType) compile(*Schema, *Element) {
}

var staticTypes = map[string]staticType{
	"string":             "string",
	"language":           "string",
	"dateTime":           "time.Time",
	"date":               "time.Time",
	"base64Binary":       "string",
	"normalizedString":   "string",
	"token":              "string",
	"NCName":             "string",
	"NMTOKENS":           "string",
	"anySimpleType":      "string",
	"anyType":            "string",
	"int":                "int",
	"integer":            "int64",
	"long":               "int64",
	"negativeInteger":    "int64",
	"nonNegativeInteger": "uint64",
	"anyURI":             "string",
	"double":             "float64",
	"decimal":            "float64", // no: http://books.xmlschemata.org/relaxng/ch19-77057.html
	"float":              "float64",
	"boolean":            "bool",
	"ID":                 "string",
	"IDREF":              "string",
	"positiveInteger":    "uint64",
	"unsignedInt":        "uint64",
	"gYear":              "string",
	"gYearMonth":         "string",
	"gMonthDay":          "string",
	"gDay":               "string",
	"gMonth":             "string",
	"time":               "time.Time",
}

func StaticType(name string) staticType {
	typ, found := staticTypes[name]
	if found {
		return typ
	}
	panic("Type xsd:" + name + " not implemented")
}

func IsStaticType(name string) bool {
	_, found := staticTypes[name]
	return found
}
