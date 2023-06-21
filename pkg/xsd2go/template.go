package xsd2go

import (
	"bytes"
	"fmt"
	tmpl "github.com/moov-io/xsd2go/pkg/templates"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/moov-io/xsd2go/pkg/xsd"
)

func GenerateTypes(templateName string, schema *xsd.Schema, outputDir string, outputFile string, tmplDir string) error {
	err := os.MkdirAll(outputDir, os.FileMode(0722))
	if err != nil {
		return err
	}
	goFile := filepath.Clean(filepath.Join(outputDir, outputFile))
	fmt.Printf("\tGenerating '%s'\n", goFile)
	f, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("Could not create '%s': %v", goFile, err)
	}

	t := template.New(templateName).Funcs(template.FuncMap{
		// Allow any template ending in suffix ".incl" to be included inline. The main template will call this
		// function at a specific point.
		"InclType": func(tmplName string, data any) (string, error) {
			tmplName = fmt.Sprintf("%s.incl", tmplName)
			t2 := template.New(tmplName)
			t2, tmplErr := parseTemplate(t2, tmplDir, tmplName)
			if tmplErr != nil {
				return "", tmplErr
			}
			var tmplBuf bytes.Buffer
			tmplErr = t2.Execute(&tmplBuf, data)
			return tmplBuf.String(), tmplErr
		},
		// Allow any template ending in suffix ".elem" to be included inline. The main template will call this
		// function at a specific point.
		"InclElem": func(tmplName string, data *xsd.Element) (string, error) {
			tmplName = fmt.Sprintf("%s.elem", tmplName)
			t2 := template.New(tmplName)
			t2, tmplErr := parseTemplate(t2, tmplDir, tmplName)
			if tmplErr != nil {
				return "", tmplErr
			}
			var tmplBuf bytes.Buffer
			tmplErr = t2.Execute(&tmplBuf, data)
			return tmplBuf.String(), tmplErr
		},
	})
	t, err = parseTemplate(t, tmplDir, templateName)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.ExecuteTemplate(&buf, templateName, schema); err != nil {
		return fmt.Errorf("Could not execute template: %v", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Could not gofmt output file\nError was: '%v'\nFile was:\n%s\n", err, buf.String())
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	_, err = exec.Command("goimports", "-w", f.Name()).Output()
	if err != nil {
		return err
	}

	return nil
}

func GetAllTemplates(tmplDir string) (map[string]xsd.Override, error) {
	dir, err := tmpl.Templates.ReadDir(tmplDir)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]xsd.Override)
	for indx := range dir {
		name, found := strings.CutSuffix(dir[indx].Name(), ".incl")
		if found {
			if _, ok := templates[name]; !ok {
				templates[name] = xsd.Override{TemplateName: name}
			}
			override := templates[name]
			override.IsIncl = true
			templates[name] = override
		}
		name, found = strings.CutSuffix(dir[indx].Name(), ".elem")
		if found {
			if _, ok := templates[name]; !ok {
				templates[name] = xsd.Override{TemplateName: name}
			}
			override := templates[name]
			override.IsElem = true
			templates[name] = override
		}
	}

	return templates, nil
}

func parseTemplateFile(tmplDir string, templateName string) (string, error) {
	in, err := tmpl.Templates.Open(tmplDir + "/" + templateName)
	if err != nil {
		return "", err
	}
	defer in.Close()

	tempText, err := io.ReadAll(in)
	return string(tempText), err
}

func parseTemplate(t *template.Template, tmplDir string, templateName string) (*template.Template, error) {
	tmplText, err := parseTemplateFile(tmplDir, templateName)
	if err != nil {
		return t, err
	}
	t, err = t.Parse(tmplText)
	if err != nil {
		return t, err
	}
	return t, nil
}
