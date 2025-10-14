package xsd

import (
	"encoding/xml"
	"fmt"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute.
type Attribute struct {
	XMLName        xml.Name    `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Name           string      `xml:"name,attr"`
	Type           reference   `xml:"type,attr"`
	Use            string      `xml:"use,attr"`
	Annotation     *Annotation `xml:"annotation"`
	DuplicateCount uint        `xml:"-"`
	Ref            reference   `xml:"ref,attr"`
	refAttr        *Attribute
	typ            Type
	schema         *Schema
}

func (a *Attribute) ContainsDocumentation() bool {
	return a.Documentation() != ""
}

func (a *Attribute) Documentation() string {
	if a.Annotation == nil {
		return ""
	}
	if len(a.Annotation.Documentations) == 0 {
		return ""
	}
	return a.Annotation.Documentations[0].GetContent()
}

// Public Go Name of this struct item.
func (a *Attribute) GoName() string {
	name := a.Name
	if a.Name == "" {
		name = a.Ref.GoName()
	}
	if a.DuplicateCount >= 2 {
		name = fmt.Sprintf("%s%d", name, a.DuplicateCount)
	}
	return strcase.ToCamel(name)
}

func (a *Attribute) GoType() string {
	if a.typ == nil {
		return "string"
	}
	return a.typ.GoName()
}

func (a *Attribute) isPlainString() bool {
	if a.typ == nil {
		return true
	}
	_, ok := a.typ.(staticType)
	return ok
}

func (a *Attribute) GoForeignModule() string {
	if a.isPlainString() {
		return ""
	}

	foreignSchema := (*Schema)(nil)
	if a.refAttr != nil {
		foreignSchema = a.refAttr.schema
	} else if a.typ != nil {
		foreignSchema = a.typ.Schema()
	}

	if foreignSchema != nil && foreignSchema != a.schema &&
		foreignSchema.TargetNamespace != a.schema.TargetNamespace {
		return foreignSchema.GoPackageName() + "."
	}
	return ""
}

func (a *Attribute) Modifiers() string {
	res := "attr"
	if a.optional() {
		res += ",omitempty"
	}
	return res
}

func (a *Attribute) XmlName() string {
	if a.Name == "" {
		return a.Ref.Name()
	}
	return a.Name
}

func (a *Attribute) optional() bool {
	// 'use' defaults to 'optional': https://www.w3.org/TR/xmlschema11-1/#declare-attribute
	return a.Use == "" || a.Use == "optional"
}

func (a *Attribute) compile(s *Schema) {
	a.schema = s
	if a.Ref != "" {
		a.refAttr = a.schema.findReferencedAttribute(a.Ref)
		if a.refAttr == nil {
			panic("Cannot resolve attribute reference: " + a.Ref)
		}
	}
	if a.Type != "" && a.typ == nil {
		a.typ = a.schema.findReferencedType(a.Type)
		if a.typ == nil {
			panic("Cannot resolve attribute type: " + a.Type)
		}
	}
}
