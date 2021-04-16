package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Enumeration struct {
	XMLName      xml.Name `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Value        string   `xml:"value,attr"`
	goNameSuffix string
}

// Public Go Name of this struct item
func (e *Enumeration) GoName() string {
	return strcase.ToCamel(e.Value + e.goNameSuffix)
}

func (e *Enumeration) Modifiers() string {
	return "-"
}

func (e *Enumeration) XmlName() string {
	return e.Value
}

func (e *Enumeration) compile(s *Schema) {
}
