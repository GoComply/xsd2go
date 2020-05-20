package xsd2go

import (
	"fmt"
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

	meta, err := xsd.Parse(f)
	if err != nil {
		return err
	}

	if err := template.GenerateTypes(meta, outputDir); err != nil {
		return err
	}
	return nil
}
