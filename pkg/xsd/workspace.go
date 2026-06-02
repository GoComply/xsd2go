package xsd

import (
	"fmt"
	"path/filepath"
)

type Workspace struct {
	Cache          map[string]*Schema // Parsed XSD schemas by its filename (user specifies initial one, and we load dependencies)
	Loaded         map[string]*Schema // Parsed AND resolved schemas by its filename
	GoModulesPath  string             // user requested go package path (example: github.com/gocomply/scap)
	xmlnsOverrides xmlnsOverrides     // user-supplied xmlns overrides
}

func NewWorkspace(goModulesPath, xsdPath string, xmlnsOverrides []string) (*Workspace, error) {
	ws := Workspace{
		Cache:         map[string]*Schema{},
		Loaded:        map[string]*Schema{},
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

// merges unique elements of newer into origin, compared by getName.
func merge[T any, M comparable](newer, origin []T, getName func(T) M) []T {
	names := make(map[M]struct{})
	for _, o := range origin {
		names[getName(o)] = struct{}{}
	}

	for _, n := range newer {
		name := getName(n)
		if _, ok := names[name]; ok {
			continue
		}
		origin = append(origin, n)
		names[name] = struct{}{}
	}
	return origin
}

func (ws *Workspace) loadXsd(xsdPath string, shouldBeInlined bool) (*Schema, error) {
	xsdPath = filepath.Clean(xsdPath)

	if schema, found := ws.Loaded[xsdPath]; found {
		return schema, nil
	}
	if cached, found := ws.Cache[xsdPath]; found {
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

	ws.Loaded[xsdPath] = schema

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
		schema.Imports = merge(isch.Imports, schema.Imports, func(i Import) string { return i.Namespace + i.SchemaLocation })
		schema.Elements = merge(isch.Elements, schema.Elements, func(e Element) string { return e.Name })
		schema.Attributes = merge(isch.Attributes, schema.Attributes, func(a Attribute) string { return a.Name })
		schema.AttributeGroups = merge(isch.AttributeGroups, schema.AttributeGroups, func(ag AttributeGroup) string { return ag.Name })
		schema.ComplexTypes = merge(isch.ComplexTypes, schema.ComplexTypes, func(ct ComplexType) string { return ct.Name })
		schema.SimpleTypes = merge(isch.SimpleTypes, schema.SimpleTypes, func(st SimpleType) string { return st.Name })
		schema.inlinedElements = merge(isch.inlinedElements, schema.inlinedElements, func(e Element) string { return e.Name })
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
