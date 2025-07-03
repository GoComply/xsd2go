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

func setXmlNameAnyForSingleElements(elements []Element) []Element {
	if len(elements) == 1 {
		result := make([]Element, 1)
		element := elements[0]
		element.XmlNameOverride = ",any"
		result[0] = element
		return result
	}
	for idx := range elements {
		element := &elements[idx]
		element.XmlNameOverride = ""
	}
	return elements
}

type ComplexType struct {
	XMLName          xml.Name     `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Name             string       `xml:"name,attr"`
	Mixed            bool         `xml:"mixed,attr"`
	AttributesDirect []Attribute  `xml:"attribute"`
	Annotation       *Annotation  `xml:"annotation"`
	Sequence         *Sequence    `xml:"sequence"`
	SequenceAll      *SequenceAll `xml:"all"`
	schema           *Schema
	SimpleContent    *SimpleContent  `xml:"simpleContent"`
	ComplexContent   *ComplexContent `xml:"complexContent"`
	Choice           *Choice         `xml:"choice"`
	content          GenericContent
}

func (ct *ComplexType) Attributes() []Attribute {
	if ct.content != nil {
		return ct.content.Attributes()
	}
	return ct.AttributesDirect
}

func (ct *ComplexType) HasXmlNameAttribute() bool {
	for _, attribute := range ct.Attributes() {
		if attribute.GoName() == "XMLName" {
			return true
		}
	}
	return false
}

func (ct *ComplexType) ContainsDocumentation() bool {
	return ct.Documentation() != ""
}

func (ct *ComplexType) Documentation() string {
	if ct.Annotation == nil {
		return ""
	}
	if len(ct.Annotation.Documentations) == 0 {
		return ""
	}
	return ct.Annotation.Documentations[0].GetContent()
}

func (ct *ComplexType) Elements() []Element {
	if ct.Sequence != nil {
		return setXmlNameAnyForSingleElements(ct.Sequence.Elements())
	} else if ct.SequenceAll != nil {
		return setXmlNameAnyForSingleElements(ct.SequenceAll.Elements())
	} else if ct.content != nil {
		return setXmlNameAnyForSingleElements(ct.content.Elements())
	} else if ct.Choice != nil {
		return ct.Choice.Elements()
	}
	return []Element{}
}

func (ct *ComplexType) GoName() string {
	return strcase.ToCamel(ct.Name)
}

func (ct *ComplexType) GoTypeName() string {
	return ct.GoName()
}

func (ct *ComplexType) ContainsInnerXml() bool {
	return ct.Mixed
}

func (ct *ComplexType) ContainsText() bool {
	return ct.content != nil && ct.content.ContainsText()
}

func (ct *ComplexType) Schema() *Schema {
	return ct.schema
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
}

type SimpleType struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Name        string       `xml:"name,attr"`
	Annotation  *Annotation  `xml:"annotation"`
	Restriction *Restriction `xml:"restriction"`
	schema      *Schema
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

func (st *SimpleType) compile(sch *Schema, parentElement *Element) {
	if st.schema == nil {
		st.schema = sch
	}

	if st.Restriction != nil {
		st.Restriction.compile(sch, parentElement)
	}
}

func (*SimpleType) Attributes() []Attribute {
	return []Attribute{}
}

func (st *SimpleType) ContainsDocumentation() bool {
	return st.Documentation() != ""
}

func (st *SimpleType) Documentation() string {
	if st.Annotation == nil {
		return ""
	}
	if len(st.Annotation.Documentations) == 0 {
		return ""
	}
	return st.Annotation.Documentations[0].GetContent()
}

func (*SimpleType) Elements() []Element {
	return []Element{}
}

func (st *SimpleType) Enums() []Enumeration {
	if st.Restriction != nil {
		return st.Restriction.Enums()
	}
	return []Enumeration{}
}

func (*SimpleType) ContainsText() bool {
	return true
}

type staticType string

func (st staticType) GoName() string {
	return string(st)
}

func (st staticType) GoTypeName() string {
	return st.GoName()
}

func (staticType) Attributes() []Attribute {
	return []Attribute{}
}

func (staticType) Elements() []Element {
	return []Element{}
}

func (staticType) Schema() *Schema {
	return nil
}

func (staticType) ContainsText() bool {
	return true
}

func (staticType) compile(*Schema, *Element) {
}

var staticTypes = map[string]staticType{
	"string":             "string",
	"language":           "string",
	"dateTime":           "string",
	"date":               "string",
	"base64Binary":       "string",
	"duration":           "string",
	"normalizedString":   "string",
	"token":              "string",
	"Name":               "string",
	"NCName":             "string",
	"NMTOKENS":           "string",
	"anySimpleType":      "string",
	"anyType":            "string",
	"int":                "int",
	"integer":            "int64",
	"long":               "int64",
	"negativeInteger":    "int64",
	"nonNegativeInteger": "uint64",
	"nonPositiveInteger": "int64",
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
	"time":               "string",
	"unsignedLong":       "uint64",
	"unsignedShort":      "uint16",
	"unsignedByte":       "uint8",
	"short":              "int16",
	"byte":               "int8",
	"hexBinary":          "string",
	"QName":              "string",
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
