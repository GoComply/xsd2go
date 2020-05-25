package xsd2go

import (
	"fmt"

	"github.com/gocomply/xsd2go/pkg/template"
	"github.com/gocomply/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string) error {
	fmt.Printf("Processing '%s'\n", xsdPath)
	schema, err := xsd.Parse(xsdPath)
	if err != nil {
		return err
	}

	for _, imp := range schema.Imports {
		if imp.Namespace == "http://www.w3.org/XML/1998/namespace" {
			continue
		}
		if err := template.GenerateTypes(imp.ImportedSchema, outputDir); err != nil {
			return err
		}
	}

	if err := template.GenerateTypes(schema, outputDir); err != nil {
		return err
	}
	return nil
}
