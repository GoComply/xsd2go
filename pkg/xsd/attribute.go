package xsd

import (
	"encoding/xml"
	"fmt"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Attribute struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Name           string   `xml:"name,attr"`
	Type           string   `xml:"type,attr"`
	Use            string   `xml:"use,attr"`
	DuplicateCount uint     `xml:"-"`
}

// Public Go Name of this struct item
func (a *Attribute) GoName() string {
	name := a.Name
	if a.DuplicateCount >= 2 {
		name = fmt.Sprintf("%s%d", name, a.DuplicateCount)
	}
	return strcase.ToCamel(name)
}
