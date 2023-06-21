package xsd2go

import (
	"fmt"

	"github.com/moov-io/xsd2go/pkg/xsd"
)

func Convert(xsdFile, outputDir, outputFile, goPackage, nsPrefix string, tmplDir string) error {
	fmt.Printf("Processing '%s'\n", xsdFile)

	templates, err := GetAllTemplates(tmplDir)
	if err != nil {
		return err
	}

	ws, err := xsd.NewWorkspace(fmt.Sprintf("%s/%s", goPackage, outputDir), xsdFile, templates)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}

		sch.NsPrefix = nsPrefix

		if err := GenerateTypes("element.tmpl", sch, outputDir, outputFile, tmplDir); err != nil {
			return err
		}
	}

	return nil
}
