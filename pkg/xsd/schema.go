package xsd

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Schema is the root XSD element
type Schema struct {
	XMLName         xml.Name  `xml:"http://www.w3.org/2001/XMLSchema schema"`
	TargetNamespace string    `xml:"targetNamespace,attr"`
	Elements        []Element `xml:"element"`
}

func Parse(r io.Reader) (*Schema, error) {
	var schema Schema
	d := xml.NewDecoder(r)

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("Error decoding XSD: %s", err)
	}

	return &schema, nil
}

func (sch *Schema) GoPackageName() string {
	return "TODO"
}
