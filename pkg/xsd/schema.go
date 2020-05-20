package xsd

import (
	"encoding/xml"
	"fmt"
	"os"
)

// Schema is the root XSD element
type Schema struct {
	XMLName         xml.Name  `xml:"http://www.w3.org/2001/XMLSchema schema"`
	TargetNamespace string    `xml:"targetNamespace,attr"`
	Imports         []Import  `xml:import`
	Elements        []Element `xml:"element"`
}

func Parse(xsdPath string) (*Schema, error) {
	f, err := os.Open(xsdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var schema Schema
	d := xml.NewDecoder(f)

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("Error decoding XSD: %s", err)
	}

	return &schema, nil
}

func (sch *Schema) GoPackageName() string {
	return "TODO"
}

func (sch *Schema) GoImportsNeeded() []string {
	return []string{"encoding/xml"}
}

type Import struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema import"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
}
