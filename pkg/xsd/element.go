package xsd

import (
	"encoding/xml"
)

// Element defines single XML element
type Element struct {
	XMLName     xml.Name     `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name        string       `xml:"name,attr"`
	Ref         reference    `xml:"ref,attr"`
	ComplexType *ComplexType `xml:"complexType"`
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
	return e.Name
}

func (e *Element) XmlName() string {
	return e.Name
}

func (e *Element) compile(s *Schema) {
	e.schema = s
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
