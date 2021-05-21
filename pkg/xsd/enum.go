package xsd

import (
	"encoding/xml"
	"strings"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Enumeration struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Value      string      `xml:"value,attr"`
	Annotation *Annotation `xml:"annotation"`
}

func (e *Enumeration) GoComments() []string {
	return e.Annotation.GoComments()
}

// Public Go Name of this struct item
func (e *Enumeration) GoName() string {
	name := e.Annotation.GetName()
	if name != "" {
		return name
	} else {
		return strcase.ToCamel(strings.ToLower(e.Value))
	}
}

func (e *Enumeration) Modifiers() string {
	return "-"
}

func (e *Enumeration) XmlName() string {
	return e.Value
}

func (e *Enumeration) compile(s *Schema) {
}
