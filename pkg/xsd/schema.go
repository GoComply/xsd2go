package xsd

import (
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// Schema is the root XSD element
type Schema struct {
	XMLName         xml.Name           `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Xmlns           Xmlns              `xml:"-"`
	TargetNamespace string             `xml:"targetNamespace,attr"`
	Imports         []Import           `xml:"import"`
	Elements        []Element          `xml:"element"`
	Attributes      []Attribute        `xml:"attribute"`
	ComplexTypes    []ComplexType      `xml:"complexType"`
	SimpleTypes     []SimpleType       `xml:"simpleType"`
	importedModules map[string]*Schema `xml:"-"`
	ModulesPath     string             `xml:"-"`
}

func parseSchema(f io.Reader) (*Schema, error) {
	schema := Schema{importedModules: map[string]*Schema{}}
	d := xml.NewDecoder(f)

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("Error decoding XSD: %s", err)
	}

	return &schema, nil
}

func (sch *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	sch.Xmlns = parseXmlns(start)

	type s Schema
	ss := (*s)(sch)
	return d.DecodeElement(ss, &start)
}

func (sch *Schema) compile() {
	for idx, _ := range sch.Elements {
		el := &sch.Elements[idx]
		el.compile(sch)
	}
	for idx, _ := range sch.ComplexTypes {
		el := &sch.ComplexTypes[idx]
		el.compile(sch)
	}
	for idx, _ := range sch.SimpleTypes {
		el := &sch.SimpleTypes[idx]
		el.compile(sch)
	}
}

func (sch *Schema) findReferencedAttribute(ref reference) *Attribute {
	innerSchema := sch.findReferencedSchemaByPrefix(ref.NsPrefix())
	if innerSchema == nil {
		panic("Internal error: referenced attribute '" + ref + "' cannot be found.")
	}
	return innerSchema.GetAttribute(ref.Name())
}

func (sch *Schema) findReferencedElement(ref reference) *Element {
	innerSchema := sch.findReferencedSchemaByPrefix(ref.NsPrefix())
	if innerSchema == nil {
		panic("Internal error: referenced element '" + ref + "' cannot be found.")
	}
	if innerSchema != sch {
		sch.registerImportedModule(innerSchema)

	}
	return innerSchema.GetElement(ref.Name())
}

func (sch *Schema) findReferencedType(ref reference) Type {
	innerSchema := sch.findReferencedSchemaByPrefix(ref.NsPrefix())
	if innerSchema == nil {
		panic("Internal error: referenced type '" + ref + "' cannot be found.")
	}
	if innerSchema != sch {
		sch.registerImportedModule(innerSchema)
	}
	return innerSchema.GetType(ref.Name())
}

func (sch *Schema) findReferencedSchemaByPrefix(xmlnsPrefix string) *Schema {
	return sch.findReferencedSchemaByXmlns(sch.xmlnsByPrefix(xmlnsPrefix))
}

func (sch *Schema) xmlnsByPrefix(xmlnsPrefix string) string {
	switch xmlnsPrefix {
	case "":
		return sch.TargetNamespace
	case "xml":
		return "http://www.w3.org/XML/1998/namespace"
	default:
		uri := sch.Xmlns.UriByPrefix(xmlnsPrefix)
		if uri == "" {
			panic("Internal error: Unknown xmlns prefix: " + xmlnsPrefix)
		}
		return uri
	}
	return ""
}

func (sch *Schema) findReferencedSchemaByXmlns(xmlns string) *Schema {
	if sch.TargetNamespace == xmlns {
		return sch
	}
	for _, imp := range sch.Imports {
		if imp.Namespace == xmlns {
			return imp.ImportedSchema
		}
	}
	return nil
}

func (sch *Schema) Empty() bool {
	return len(sch.Elements) == 0 && len(sch.ComplexTypes) == 0
}

func (sch *Schema) GetAttribute(name string) *Attribute {
	for idx, attr := range sch.Attributes {
		if attr.Name == name {
			return &sch.Attributes[idx]
		}
	}
	return nil
}

func (sch *Schema) GetElement(name string) *Element {
	for idx, elm := range sch.Elements {
		if elm.Name == name {
			return &sch.Elements[idx]
		}
	}
	return nil
}

func (sch *Schema) GetType(name string) Type {
	if name == "string" {
		return staticType("string")
	}
	for idx, typ := range sch.ComplexTypes {
		if typ.Name == name {
			return &sch.ComplexTypes[idx]
		}
	}
	for idx, typ := range sch.SimpleTypes {
		if typ.Name == name {
			return &sch.SimpleTypes[idx]
		}
	}
	return nil
}

func (sch *Schema) GoPackageName() string {
	xmlnsPrefix := sch.Xmlns.PrefixByUri(sch.TargetNamespace)
	return strings.ReplaceAll(xmlnsPrefix, "-", "_")
}

func (sch *Schema) GoImportsNeeded() []string {
	imports := []string{"encoding/xml"}
	for _, importedMod := range sch.importedModules {
		imports = append(imports, fmt.Sprintf("%s/%s", sch.ModulesPath, importedMod.GoPackageName()))
	}
	sort.Strings(imports)
	return imports
}

func (sch *Schema) registerImportedModule(module *Schema) {
	sch.importedModules[module.GoPackageName()] = module
}

type Import struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema import"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	ImportedSchema *Schema  `xml:"-"`
}

func (i *Import) load(ws *Workspace, baseDir string) (err error) {
	i.ImportedSchema, err = ws.loadXsd(filepath.Join(baseDir, i.SchemaLocation))
	return
}
