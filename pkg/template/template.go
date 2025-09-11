package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gocomply/xsd2go/pkg/xsd"
)

//go:embed types.tmpl
var templText string

func GenerateTypes(schema *xsd.Schema, outputDir string) error {
	t, err := newTemplate()
	if err != nil {
		return err
	}

	packageName := schema.GoPackageName()
	dir := filepath.Join(outputDir, packageName)
	err = os.MkdirAll(dir, os.FileMode(0722))
	if err != nil {
		return err
	}
	goFile := filepath.Clean(filepath.Join(dir, "models.go"))
	fmt.Printf("\tGenerating '%s'\n", goFile)
	f, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("could not create '%s': %w", goFile, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := t.Execute(&buf, schema); err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("unable to gofmt output file %s, error: %w", buf.String(), err)
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func newTemplate() (*template.Template, error) {
	return template.New("types.tmpl").Funcs(template.FuncMap{}).Parse(templText)
}
