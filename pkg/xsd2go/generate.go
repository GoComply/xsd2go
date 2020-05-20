package xsd2go

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/gocomply/xsd2go/pkg/template"
	"github.com/gocomply/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string) error {
	fmt.Println("Processing ", xsdPath)
	f, err := os.Open(xsdPath)
	if err != nil {
		return err
	}
	defer f.Close()

	meta, err := decode(f)
	if err != nil {
		return err
	}

	if err := template.GenerateTypes(meta, outputDir); err != nil {
		return err
	}
	return nil
}

func decode(r io.Reader) (*xsd.Schema, error) {
	var schema xsd.Schema

	d := xml.NewDecoder(r)

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("Error decoding XSD: %s", err)
	}

	return &schema, nil
}
