package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

// Element defines single XML element
type Element struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name        string       `xml:"name,attr"`
	Type        reference    `xml:"type,attr"`
	Ref         reference    `xml:"ref,attr"`
	refElm      *Element     `xml:"-"`
	ComplexType *ComplexType `xml:"complexType"`
	refType     Type         `xml:"-"`
	schema      *Schema      `xml:"-"`
}

func (e *Element) Attributes() []Attribute {
	if e.ComplexType != nil {
		return e.ComplexType.Attributes
	}
	return []Attribute{}
}

func (e *Element) Elements() []Element {
	if e.ComplexType != nil {
		return e.ComplexType.Elements()
	}
	return []Element{}
}

func (e *Element) GoName() string {
	name := e.Name
	if name == "" {
		return e.refElm.GoName()
	}
	return strcase.ToCamel(name)
}

func (e *Element) GoTypeName() string {
	if e.Type != "" {
		return strcase.ToCamel(e.Type.Name())
	}
	return e.GoName()
}

func (e *Element) GoForeignModule() string {
	if e.refElm != nil && e.schema != e.refElm.schema {
		return e.refElm.schema.GoPackageName() + "."
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

func (e *Element) compile(s *Schema) {
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
	if e.ComplexType != nil {
		// Handle improbable name clash. Consider XSD defining two attributes on the element:
		// "id" and "Id", this would create name clash given the camelization we do.
		goNames := map[string]uint{}
		for idx, _ := range e.ComplexType.Attributes {
			attribute := &e.ComplexType.Attributes[idx]
			attribute.compile(s)

			count := goNames[attribute.GoName()]
			count += 1
			goNames[attribute.GoName()] = count
			attribute.DuplicateCount = count
			// Second GoName may be different depending on the DuplicateCount
			goNames[attribute.GoName()] = count
		}

		elements := e.Elements()
		for idx, _ := range elements {
			el := &elements[idx]
			el.compile(s)
		}
	}
}
