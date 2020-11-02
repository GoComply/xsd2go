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
	schema           *Schema `xml:"-"`
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

func (att *AttributeGroup) Elements() []Element {
	return []Element{}
}

func (att *AttributeGroup) ContainsText() bool {
	return true
}
