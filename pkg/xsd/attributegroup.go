package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

type AttributeGroup struct {
	XMLName          xml.Name    `xml:"http://www.w3.org/2001/XMLSchema attributeGroup"`
	Name             string      `xml:"name,attr"`
	Ref              reference   `xml:"ref,attr"`
	AttributesDirect []Attribute `xml:"attribute"`
	typ              Type
	schema           *Schema
}

func (att *AttributeGroup) Attributes() []Attribute {
	attrs := att.AttributesDirect
	if att.typ != nil {
		attrs = append(attrs, att.typ.Attributes()...)
	}
	return attrs
}

func (att *AttributeGroup) compile(sch *Schema, parentElement *Element) {
	att.schema = sch
	if att.Ref != "" {
		att.typ = sch.findReferencedType(att.Ref)
		if att.typ == nil {
			panic("Cannot build xsd:attributeGroup: unknown type: " + string(att.Ref))
		}
		att.typ.compile(sch, parentElement)
	}

	// Handle improbable name clash. Consider XSD defining two attributes on the element:
	// "id" and "Id", this would create name clash given the camelization we do.
	goNames := map[string]uint{}
	for idx := range att.Attributes() {
		attribute := &att.Attributes()[idx]
		attribute.compile(sch)

		count := goNames[attribute.GoName()]
		count += 1
		goNames[attribute.GoName()] = count
		attribute.DuplicateCount = count
		// Second GoName may be different depending on the DuplicateCount
		goNames[attribute.GoName()] = count
	}

}

func (att *AttributeGroup) GoName() string {
	return strcase.ToCamel(att.Name)

}

func (att *AttributeGroup) GoTypeName() string {
	return att.GoName()

}

func (att *AttributeGroup) Schema() *Schema {
	return att.schema
}

func (*AttributeGroup) Elements() []Element {
	return []Element{}
}

func (*AttributeGroup) ContainsText() bool {
	return true
}
