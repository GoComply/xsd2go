package xsd

import (
	"encoding/xml"
)

// Schema is the root XSD element
type Schema struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2001/XMLSchema schema"`
	TargetNamespace string   `xml:"targetNamespace,attr"`
}

func (sch *Schema) GoPackageName() string {
	return "TODO"
}
