package xsd

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
)

// Schema is the root XSD element.
type Schema struct {
	XMLName               xml.Name          `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Xmlns                 Xmlns             `xml:"-"`
	TargetNamespace       string            `xml:"targetNamespace,attr"`
	Annotation            *Annotation       `xml:"annotation"`
	Includes              []*Include        `xml:"include"`
	Imports               []*Import         `xml:"import"`
	Elements              []*Element        `xml:"element"`
	Attributes            []*Attribute      `xml:"attribute"`
	AttributeGroups       []*AttributeGroup `xml:"attributeGroup"`
	ComplexTypes          []*ComplexType    `xml:"complexType"`
	SimpleTypes           []*SimpleType     `xml:"simpleType"`
	importedModules       map[string]*Schema
	ModulesPath           string `xml:"-"`
	filePath              string
	inlinedElements       []*Element
	goPackageNameOverride string
	typeOverrides         typeOverrides
}

// schemaProvider is a helper interface detecting types that can provide a schema.
type schemaProvider interface {
	Schema() *Schema
}

func ReadSchemaFromFile(xsdPath string) (*Schema, error) {
	xsdPathClean := filepath.Clean(xsdPath)
	f, err := os.Open(xsdPathClean)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while closing file %s, %v", xsdPathClean, err)
		}
	}()

	schema, err := parseSchema(f)
	if err != nil {
		return nil, fmt.Errorf("%w; while processing %s", err, xsdPath)
	}

	return schema, nil
}

func parseSchema(f io.Reader) (*Schema, error) {
	schema := Schema{importedModules: map[string]*Schema{}}
	d := xml.NewDecoder(f)
	d.CharsetReader = charset.NewReaderLabel

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("error decoding XSD: %w", err)
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
	if sch.TargetNamespace == "" {
		fmt.Fprintf(os.Stderr, "Warning: missing explicit /xsd:chema/@targetNamespace; using '%s' instead\n", sch.GoPackageName())
		sch.TargetNamespace = sch.GoPackageName()
	}

	for idx := range sch.Elements {
		el := sch.Elements[idx]
		el.compile(sch, nil)
	}
	for idx := range sch.AttributeGroups {
		att := sch.AttributeGroups[idx]
		att.compile(sch, nil)
	}
	for idx := range sch.ComplexTypes {
		ct := sch.ComplexTypes[idx]
		ct.compile(sch, nil)
	}
	for idx := range sch.SimpleTypes {
		st := sch.SimpleTypes[idx]
		st.compile(sch, nil)
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
		panic("Internal error: referenced element '" + string(ref) + "' cannot be found.")
	}
	if innerSchema != sch {
		sch.registerImportedModule(innerSchema)
	}
	return innerSchema.GetElement(ref.Name())
}

func (sch *Schema) findReferencedType(ref reference) Type {
	innerSchema := sch.findReferencedSchemaByPrefix(ref.NsPrefix())
	if innerSchema == nil {
		xmlnsUri := sch.Xmlns.UriByPrefix(ref.NsPrefix())
		if typeName, overridden := sch.typeOverrides.overrideType(xmlnsUri, ref.Name()); overridden {
			return staticType(typeName)
		}
		if xmlnsUri == "http://www.w3.org/2001/XMLSchema" { //nolint:revive
			return StaticType(ref.Name())
		}
		panic("Internal error: referenced type '" + string(ref) + "' cannot be found.")
	}
	if innerSchema != sch {
		sch.registerImportedModule(innerSchema)
	}
	if typeName, overridden := sch.typeOverrides.overrideType(innerSchema.TargetNamespace, ref.Name()); overridden {
		return staticType(typeName)
	}
	return innerSchema.GetType(ref.Name())
}

func (sch *Schema) findReferencedSchemaByPrefix(xmlnsPrefix string) *Schema {
	return sch.findReferencedSchemaByXmlns(sch.xmlnsByPrefix(xmlnsPrefix))
}

func (sch *Schema) xmlnsByPrefix(xmlnsPrefix string) string {
	uri := sch.xmlnsByPrefixInternal(xmlnsPrefix)
	if uri == "" {
		panic("Internal error: Unknown xmlns prefix: " + xmlnsPrefix)
	}
	return uri
}

func (sch *Schema) xmlnsByPrefixInternal(xmlnsPrefix string) string {
	switch xmlnsPrefix {
	case "":
		return sch.TargetNamespace
	case "xml":
		return "http://www.w3.org/XML/1998/namespace" //nolint:revive
	default:
		uri := sch.Xmlns.UriByPrefix(xmlnsPrefix)
		if uri == "" {
			for _, imported := range sch.importedModules {
				uri = imported.xmlnsByPrefixInternal(xmlnsPrefix)
				if uri != "" {
					return uri
				}
			}
		}
		return uri
	}
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
	for _, imp := range sch.importedModules {
		s := imp.findReferencedSchemaByXmlns(xmlns)
		if s != nil {
			return s
		}
	}
	return nil
}

func (sch *Schema) Empty() bool {
	return len(sch.Elements) == 0 && len(sch.ComplexTypes) == 0 && len(sch.ExportableSimpleTypes()) == 0
}

func (sch *Schema) encodingXmlImportNeeded() bool {
	return len(sch.Elements) != 0 || len(sch.ComplexTypes) != 0
}

func deduplicateElementsLossfree(elements []*Element) []*Element {
	seen := make(map[string]int, len(elements))
	for j, element := range elements {
		dupeCount, dupe := seen[element.GoName()]
		if !dupe {
			seen[element.GoName()] = 1
			continue
		}
		elements[j].nameOverride = fmt.Sprintf("%s-%d", element.GoName(), dupeCount)
		seen[element.GoName()]++
	}
	return elements
}

func (sch *Schema) ExportableElements() []*Element {
	return deduplicateElementsLossfree(append(sch.Elements, sch.inlinedElements...))
}

func (sch *Schema) ExportableComplexTypes() []*ComplexType {
	elCache := map[string]bool{}
	for _, el := range sch.Elements {
		elCache[el.GoName()] = true
	}

	var res []*ComplexType
	for _, typ := range sch.ComplexTypes {
		_, found := elCache[typ.GoName()]
		if !found {
			res = append(res, typ)
		}
	}
	return res
}

func (sch *Schema) ExportableSimpleTypes() []*SimpleType {
	elCache := map[string]bool{}
	for _, el := range sch.Elements {
		elCache[el.GoName()] = true
	}

	var res []*SimpleType
	for _, typ := range sch.SimpleTypes {
		_, found := elCache[typ.GoName()]
		if !found {
			res = append(res, typ)
		}
	}
	return res
}

func (sch *Schema) GetAttribute(name string) *Attribute {
	for idx, attr := range sch.Attributes {
		if attr.Name == name {
			return sch.Attributes[idx]
		}
	}
	return nil
}

func (sch *Schema) GetElement(name string) *Element {
	for idx, elm := range sch.Elements {
		if elm.Name == name {
			return sch.Elements[idx]
		}
	}
	return nil
}

func (sch *Schema) GetType(name string) Type {
	for idx, typ := range sch.ComplexTypes {
		if typ.Name == name {
			return sch.ComplexTypes[idx]
		}
	}
	for idx, typ := range sch.SimpleTypes {
		if typ.Name == name {
			return sch.SimpleTypes[idx]
		}
	}
	for idx, typ := range sch.AttributeGroups {
		if typ.Name == name {
			return sch.AttributeGroups[idx]
		}
	}
	if IsStaticType(name) {
		return StaticType(name)
	}
	return nil
}

func (sch *Schema) ContainsDocumentation() bool {
	return sch.Documentation() != ""
}

func (sch *Schema) Documentation() string {
	if sch.Annotation == nil {
		return ""
	}
	if len(sch.Annotation.Documentations) == 0 {
		return ""
	}
	return sch.Annotation.Documentations[0].GetContent()
}

func (sch *Schema) GoPackageName() string {
	if sch.goPackageNameOverride != "" {
		return sch.goPackageNameOverride
	}
	xmlnsPrefix := sch.Xmlns.PrefixByUri(sch.TargetNamespace)
	if xmlnsPrefix == "" {
		xmlnsPrefix = strings.TrimSuffix(filepath.Base(sch.filePath), ".xsd")
	}
	return strings.ReplaceAll(strings.ReplaceAll(xmlnsPrefix, "-", "_"), ".", "_")
}

func (sch *Schema) GoImportsNeeded() []string {
	imports := []string{}
	if sch.encodingXmlImportNeeded() {
		imports = append(imports, "encoding/xml")
	}
	for _, importedMod := range sch.importedModules {
		imports = append(imports, fmt.Sprintf("%s/%s", sch.ModulesPath, importedMod.GoPackageName()))
	}
	sort.Strings(imports)
	return imports
}

func (sch *Schema) registerImportedModule(module *Schema) {
	sch.importedModules[module.GoPackageName()] = module
}

func (sch *Schema) merge(other *Schema) {
	sch.Imports = append(other.Imports, sch.Imports...)
	sch.Elements = append(other.Elements, sch.Elements...)
	sch.Attributes = append(other.Attributes, sch.Attributes...)
	sch.AttributeGroups = append(other.AttributeGroups, sch.AttributeGroups...)
	sch.ComplexTypes = append(other.ComplexTypes, sch.ComplexTypes...)
	sch.SimpleTypes = append(other.SimpleTypes, sch.SimpleTypes...)
	sch.inlinedElements = append(other.inlinedElements, sch.inlinedElements...)
	for key, sch := range other.importedModules {
		sch.importedModules[key] = sch
	}
}

// Some elements are not defined at the top-level, rather these are inlined in the complexType definitions.
func (sch *Schema) registerInlinedElement(el *Element, parentElement *Element) {
	if sch.isElementInlined(el) {
		if el.Name == "" {
			panic("Not implemented: found inlined xsd:element without @name attribute")
		}
		el.prefixNameWithParent(parentElement)

		cloned := *el
		sch.inlinedElements = append(sch.inlinedElements, &cloned)
	}
}

func (sch *Schema) isElementInlined(el *Element) bool {
	found := false
	for idx := range sch.Elements {
		e := sch.Elements[idx]
		if *e == *el {
			found = true
			break
		}
	}
	return !found
}

type Import struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema import"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	ImportedSchema *Schema  `xml:"-"`
}

func (i *Import) load(ws *Workspace, baseDir string) (err error) {
	if i.SchemaLocation != "" {
		fname := filepath.ToSlash(filepath.Join(baseDir, i.SchemaLocation))
		i.ImportedSchema, err = ws.loadXsd(fname, false)
	}
	return
}

type Include struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema include"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	IncludedSchema *Schema  `xml:"-"`
}

func (i *Include) load(ws *Workspace, baseDir string) (err error) {
	if i.SchemaLocation != "" {
		fname := filepath.ToSlash(filepath.Join(baseDir, i.SchemaLocation))
		i.IncludedSchema, err = ws.loadXsd(fname, true)
	}
	return
}
