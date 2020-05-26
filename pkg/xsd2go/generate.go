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
	schema.ModulesPath = fmt.Sprintf("%s/%s", goModule, outputDir)

	for _, imp := range schema.Imports {
		if imp.ImportedSchema.Empty() {
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
