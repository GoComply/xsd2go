package xsd

import (
	"fmt"
	"path/filepath"
)

type Workspace struct {
	Cache          map[string]*Schema // Parsed XSD schemas by its filename (user specifies initial one, and we load dependencies)
	GoModulesPath  string             // user requested go package path (example: github.com/gocomply/scap)
	xmlnsOverrides xmlnsOverrides     // user-supplied xmlns overrides
}

func NewWorkspace(goModulesPath, xsdPath string, xmlnsOverrides []string) (*Workspace, error) {
	ws := Workspace{
		Cache:         map[string]*Schema{},
		GoModulesPath: goModulesPath,
	}
	var err error
	ws.xmlnsOverrides, err = ParseXmlnsOverrides(xmlnsOverrides)
	if err != nil {
		return nil, err
	}

	_, err = ws.loadXsd(xsdPath, false)
	if err != nil {
		return nil, err
	}
	return &ws, ws.compile()
}

func (ws *Workspace) loadXsd(xsdPath string, shouldBeInlined bool) (*Schema, error) {
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	schema, err := ReadSchemaFromFile(xsdPath)
	if err != nil {
		return nil, err
	}

	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	schema.goPackageNameOverride = ws.xmlnsOverrides.override(schema.TargetNamespace)

	if !shouldBeInlined {
		// Cache all loaded schemas in the workspace, unless it was brought in by xsd:include element.
		// Unlike xsd:import, xsd:include does not result in a separate schema in the workspace.
		ws.Cache[xsdPath] = schema
	}

	dir := filepath.Dir(xsdPath)

	for idx := range schema.Includes {
		si := schema.Includes[idx]
		if err := si.load(ws, dir); err != nil {
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
		if err := schema.Imports[idx].load(ws, dir); err != nil {
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
		prevXmlns, dupeFound := uniqPkgNames[goPackageName]
		if dupeFound {
			return fmt.Errorf("malformed workspace; multiple XSD files refer to itself with xmlns shorthand: '%s':\n - %s\n - %s\nWhile this is valid in XSD it is impractical for golang code generation.\nConsider providing --xmlns-override=%s=mygopackage", goPackageName, prevXmlns, schema.TargetNamespace, schema.TargetNamespace)
		}
		uniqPkgNames[goPackageName] = schema.TargetNamespace
	}

	return nil
}
