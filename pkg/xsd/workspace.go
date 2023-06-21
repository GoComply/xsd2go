package xsd

import (
	"fmt"
	"os"
	"path/filepath"
)

type Workspace struct {
	Cache         map[string]*Schema
	GoModulesPath string
}

func NewWorkspace(goModulesPath, xsdFile string, templates map[string]Override) (*Workspace, error) {
	ws := Workspace{
		Cache:         map[string]*Schema{},
		GoModulesPath: goModulesPath,
	}
	var err error
	if err != nil {
		return nil, err
	}

	_, err = ws.loadXsd(xsdFile, templates, true)
	if err != nil {
		return nil, err
	}
	return &ws, ws.compile()
}

func (ws *Workspace) loadXsd(xsdPath string, templates map[string]Override, cache bool) (*Schema, error) {
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	xsdPathClean := filepath.Clean(xsdPath)
	f, err := os.Open(xsdPathClean)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	schema, err := parseSchema(f)
	if err != nil {
		return nil, err
	}

	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	schema.TemplateOverrides = templates
	// Won't cache included schemas - we need to append contents to the current schema.
	if cache {
		ws.Cache[xsdPath] = schema
	}

	dir := filepath.Dir(xsdPath)

	for idx := range schema.Includes {
		si := schema.Includes[idx]
		if err := si.load(ws, schema.TemplateOverrides, dir); err != nil {
			return nil, err
		}

		isch := si.IncludedSchema
		schema.Imports = append(isch.Imports, schema.Imports...)
		schema.Elements = append(isch.Elements, schema.Elements...)
		schema.Attributes = append(isch.Attributes, schema.Attributes...)
		schema.AttributeGroups = append(isch.AttributeGroups, schema.AttributeGroups...)
		schema.ComplexTypes = append(isch.ComplexTypes, schema.ComplexTypes...)
		schema.SimpleTypes = append(isch.SimpleTypes, schema.SimpleTypes...)
		schema.inlinedElements = append(isch.inlinedElements, schema.inlinedElements...)
		for key, sch := range isch.importedModules {
			schema.importedModules[key] = sch
		}
	}

	for idx := range schema.Imports {
		if err := schema.Imports[idx].load(ws, schema.TemplateOverrides, dir); err != nil {
			return nil, err
		}
	}
	schema.compile()
	return schema, nil
}

func (ws *Workspace) compile() error {
	uniqPkgNames := map[string]string{}

	for _, schema := range ws.Cache {
		goPackageName := schema.GoPackageName()
		prevXmlns, ok := uniqPkgNames[goPackageName]
		if ok {
			return fmt.Errorf("Malformed workspace. Multiple XSD files refer to itself with xmlns shorthand: '%s':\n - %s\n - %s\nWhile this is xsd in XSD it is impractical for golang code generation.\nConsider providing --xmlns-override=%s=mygopackage", goPackageName, prevXmlns, schema.TargetNamespace, schema.TargetNamespace)
		}
		uniqPkgNames[goPackageName] = schema.TargetNamespace
	}

	return nil
}
