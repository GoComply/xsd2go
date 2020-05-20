package template

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gocomply/xsd2go/pkg/xsd"
	"github.com/markbates/pkger"
)

func GenerateTypes(schema *xsd.Schema, outputDir string) error {
	t, err := newTemplate(outputDir)
	if err != nil {
		return err
	}

	packageName := schema.GoPackageName()
	dir := filepath.Join(outputDir, packageName)
	err = os.MkdirAll(dir, os.FileMode(0722))
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s.go", dir, packageName))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, schema); err != nil {
		return err
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.New(err.Error() + " in following file:\n" + string(buf.Bytes()))
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func newTemplate(outputDir string) (*template.Template, error) {
	in, err := pkger.Open("/pkg/template/types.tmpl")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	tempText, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	return template.New("types.tmpl").Funcs(template.FuncMap{}).Parse(string(tempText))
}
