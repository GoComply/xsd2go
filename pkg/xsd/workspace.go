package xsd

import (
	"fmt"
	"path/filepath"
)

type merge struct {
	a string
	b string
}

func (m merge) reverse() merge {
	return merge{a: m.b, b: m.a}
}

type Workspace struct {
	Cache             map[string]*Schema // Parsed XSD schemas by its filename (user specifies initial one, and we load dependencies)
	Lookup            map[string]*Schema
	GoModulesPath     string         // user requested go package path (example: github.com/gocomply/scap)
	xmlnsOverrides    xmlnsOverrides // user-supplied xmlns overrides
	typeOverrides     typeOverrides
	uncompiledSchemas map[string]*Schema
	alias             map[string]*Schema
	merges            map[merge]struct{}
}

func NewWorkspace(goModulesPath string, xsdPaths []string, xmlnsOverridesInput []string, typeOverridesInput []string) (*Workspace, error) {
	var err error
	ws := Workspace{
		Cache:             map[string]*Schema{},
		Lookup:            map[string]*Schema{},
		GoModulesPath:     goModulesPath,
		typeOverrides:     typeOverrides{},
		uncompiledSchemas: map[string]*Schema{},
		alias:             map[string]*Schema{},
		merges:            map[merge]struct{}{},
	}

	ws.xmlnsOverrides, err = ParseXmlnsOverrides(xmlnsOverridesInput)
	if err != nil {
		return nil, err
	}

	ws.typeOverrides, err = parseTypeOverrides(typeOverridesInput)
	if err != nil {
		return nil, err
	}

	for _, input := range xsdPaths {
		err = ws.addXsd(input)
		if err != nil {
			return nil, err
		}
	}
	return &ws, ws.compile()
}

func (ws *Workspace) addXsd(xsdPath string) error {
	_, err := ws.loadXsd(xsdPath, false)
	if err != nil {
		return err
	}
	return nil
}

func (ws *Workspace) loadXsd(xsdPath string, shouldBeInlined bool) (*Schema, error) {
	if alias, ok := ws.alias[xsdPath]; ok {
		return alias, nil
	}

	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	schema, err := ReadSchemaFromFile(xsdPath)
	if err != nil {
		return nil, err
	}

	// Some schemas continue in separate files.
	seenSchema, isExternalContinuation := ws.Lookup[schema.TargetNamespace]
	if !isExternalContinuation {
		// Add the schema to the lookup map so it can be found by URI more easily.
		ws.Lookup[schema.TargetNamespace] = schema
	}

	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	schema.goPackageNameOverride = ws.xmlnsOverrides.override(schema.TargetNamespace)

	if !shouldBeInlined && !isExternalContinuation {
		// Cache all loaded schemas in the workspace, unless it was brought in by xsd:include element.
		// Unlike xsd:import, xsd:include does not result in a separate schema in the workspace.
		ws.Cache[xsdPath] = schema
	}

	dir := filepath.Dir(xsdPath)

	// Merge included schemas.
	for idx := range schema.Includes {
		si := schema.Includes[idx]
		if err := si.load(ws, dir); err != nil {
			return nil, err
		}

		included := si.IncludedSchema

		// Skip this include if the included schema is an external continuation of a seen schema.
		if isExternalContinuation && included.TargetNamespace == seenSchema.TargetNamespace {
			// Effectively a self-include, so skip it.
			continue
		}

		if schema.TargetNamespace == included.TargetNamespace {
			// Prevent self-include.
			continue
		}

		// Merge this schema into the current one.
		if !ws.alreadyMerged(schema, included) {
			ws.merge(schema, included)
		}
	}

	for idx := range schema.Imports {
		if err := schema.Imports[idx].load(ws, dir); err != nil {
			return nil, err
		}
	}

	if isExternalContinuation {
		// Merge the incoming schema into the on seen earlier.
		if !ws.alreadyMerged(seenSchema, schema) {
			ws.merge(seenSchema, schema)
		}

		ws.uncompiledSchemas[seenSchema.TargetNamespace] = seenSchema

		// Treat this schema as the one seen earlier going forward.
		ws.alias[xsdPath] = seenSchema

		return seenSchema, nil
	}

	schema.compile()
	return schema, nil
}

func (ws *Workspace) compile() error {
	// Compile all uncompiled schemas
	for _, schema := range ws.uncompiledSchemas {
		schema.compile()
	}

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

func (ws *Workspace) alreadyMerged(schema *Schema, other *Schema) bool {
	m := merge{a: schema.filePath, b: other.filePath}
	if _, ok := ws.merges[m]; ok {
		return true
	} else if _, ok = ws.merges[m.reverse()]; ok {
		return true
	}
	return false
}

func (ws *Workspace) merge(schema *Schema, other *Schema) {
	schema.merge(other)

	// Track this merge.
	m := merge{a: schema.filePath, b: other.filePath}
	ws.merges[m] = struct{}{}
}
