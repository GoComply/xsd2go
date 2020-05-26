package xsd

import (
	"fmt"
	"os"
	"path/filepath"
)

type Workspace struct {
	PrimarySchema *Schema
	Cache         map[string]*Schema
}

func NewWorkspace(xsdPath string) (*Workspace, error) {
	ws := Workspace{Cache: map[string]*Schema{}}
	var err error
	ws.PrimarySchema, err = ws.loadXsd(xsdPath)
	return &ws, err
}

func (ws *Workspace) loadXsd(xsdPath string) (*Schema, error) {
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	f, err := os.Open(xsdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	schema, err := parseSchema(f)
	if err != nil {
		return nil, err
	}
	ws.Cache[xsdPath] = schema

	dir := filepath.Dir(xsdPath)
	for idx, _ := range schema.Imports {
		if err := schema.Imports[idx].load(ws, dir); err != nil {
			return nil, err
		}
	}
	schema.compile()
	return schema, nil
}
