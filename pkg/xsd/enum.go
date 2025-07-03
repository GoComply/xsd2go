package xsd

import (
	"encoding/xml"
	"strings"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute.
type Enumeration struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Value   string   `xml:"value,attr"`
}

// Public Go Name of this struct item.
func (e *Enumeration) GoName() string {
	return strcase.ToCamel(strings.ToLower(e.Value))
}

func (e *Enumeration) Modifiers() string {
	return "-"
}

func (e *Enumeration) XmlName() string {
	return e.Value
}
