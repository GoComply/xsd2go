package xsd

import (
	"encoding/xml"
	"fmt"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Attribute struct {
	XMLName        xml.Name   `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Name           string     `xml:"name,attr"`
	Type           string     `xml:"type,attr"`
	Use            string     `xml:"use,attr"`
	DuplicateCount uint       `xml:"-"`
	Ref            reference  `xml:"ref,attr"`
	refAttr        *Attribute `xml:"-"`
	schema         *Schema    `xml:"-"`
}

// Public Go Name of this struct item
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

func (a *Attribute) XmlName() string {
	return a.Name
}

func (a *Attribute) compile(s *Schema) {
	a.schema = s
	if a.Ref != "" {
		a.refAttr = a.schema.findReferencedAttribute(a.Ref)
		if a.refAttr == nil {
			panic("Cannot resolve attribute reference: " + a.Ref)
		}
	}
}
