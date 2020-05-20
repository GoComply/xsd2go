package xsd2go

import (
	"fmt"

	"github.com/gocomply/xsd2go/pkg/template"
	"github.com/gocomply/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string) error {
	fmt.Printf("Processing '%s'\n", xsdPath)
	meta, err := xsd.Parse(xsdPath)
	if err != nil {
		return err
	}

	if err := template.GenerateTypes(meta, outputDir); err != nil {
		return err
	}
	return nil
}
