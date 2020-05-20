package xsd

import (
	"encoding/xml"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Attribute struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Use     string   `xml:"use,attr"`
}

// Public Go Name of this struct item
func (a *Attribute) GoName() string {
	return strcase.ToCamel(a.Name)
}
