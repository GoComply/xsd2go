package xsd

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/iancoleman/strcase"
)

// Element defines single XML element
type Element struct {
	XMLName      xml.Name     `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name         string       `xml:"name,attr"`
	nameOverride string       `xml:"-"`
	Type         reference    `xml:"type,attr"`
	Ref          reference    `xml:"ref,attr"`
	MinOccurs    string       `xml:"minOccurs,attr"`
	MaxOccurs    string       `xml:"maxOccurs,attr"`
	refElm       *Element     `xml:"-"`
	ComplexType  *ComplexType `xml:"complexType"`
	SimpleType   *SimpleType  `xml:"simpleType"`
	refType      Type         `xml:"-"`
	schema       *Schema      `xml:"-"`
}

func (e *Element) Attributes() []Attribute {
	if e.ComplexType != nil {
		return e.ComplexType.Attributes()
	}
	return []Attribute{}
}

func (e *Element) Elements() []Element {
	if e.ComplexType != nil {
		return e.ComplexType.Elements()
	} else if e.refType != nil {
		return e.refType.Elements()
	}
	return []Element{}
}

func (e *Element) GoFieldName() string {
	name := e.Name
	if name == "" {
		return e.refElm.GoName()
	}
	return strcase.ToCamel(name)
}

func (e *Element) GoName() string {
	if e.nameOverride != "" {
		return strcase.ToCamel(e.nameOverride)
	}
	return e.GoFieldName()

}

func (e *Element) GoMemLayout() string {
	if e.isArray() {
		return "[]"
	}
	if (e.MaxOccurs == "1" || e.MaxOccurs == "") && e.MinOccurs == "0" && e.GoTypeName() != "string" {
		return "*"
	}
	return ""
}

func (e *Element) GoTypeName() string {
	if e.Type != "" {
		return e.refType.GoTypeName()
	} else if e.isPlainString() {
		return "string"
	}
	return e.GoName()
}

func (e *Element) GoForeignModule() string {
	foreignSchema := (*Schema)(nil)
	if e.refElm != nil {
		foreignSchema = e.refElm.schema
	} else if e.refType != nil {
		foreignSchema = e.refType.Schema()
	}

	if foreignSchema != nil && foreignSchema != e.schema {
		return foreignSchema.GoPackageName() + "."
	}
	return ""
}

func (e *Element) XmlName() string {
	name := e.Name
	if name == "" {
		return e.refElm.XmlName()
	}
	return name
}

func (e *Element) isPlainString() bool {
	return e.SimpleType != nil || (e.Type == "" && e.Ref == "" && e.ComplexType == nil)
}

func (e *Element) isArray() bool {
	if e.MaxOccurs == "unbounded" {
		return true
	}
	occurs, err := strconv.Atoi(e.MaxOccurs)
	return err == nil && occurs > 1
}

func (e *Element) compile(s *Schema, parentElement *Element) {
	e.schema = s
	if e.Ref != "" {
		e.refElm = e.schema.findReferencedElement(e.Ref)
		if e.refElm == nil {
			panic("Cannot resolve element reference: " + e.Ref)
		}
	}
	if e.Type != "" {
		e.refType = e.schema.findReferencedType(e.Type)
		if e.refType == nil {
			panic("Cannot resolve type reference: " + string(e.Type))
		}
	}

	if e.ComplexType != nil && e.refType != nil {
		panic("Not implemented: element " + e.GoName() + "defines complexType and refType")
	}

	if e.ComplexType != nil {
		e.ComplexType.compile(s, e)
	}
	if e.Ref == "" && e.Type == "" && !e.isPlainString() {
		e.schema.registerInlinedElement(e, parentElement)
	}
}

func (e *Element) prefixNameWithParent(parentElement *Element) {
	// In case there are inlined xsd:elements within another xsd:elements, it may happen that two top-level xsd:elements
	// define child xsd:element of a same name. In such case, we need to override children name to avoid name clashes.
	if parentElement != nil {
		e.nameOverride = fmt.Sprintf("%s-%s", parentElement.GoName(), e.GoName())
	}
}
