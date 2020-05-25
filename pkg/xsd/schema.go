package xsd

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// Schema is the root XSD element
type Schema struct {
	XMLName         xml.Name      `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Xmlns           Xmlns         `xml:"-"`
	TargetNamespace string        `xml:"targetNamespace,attr"`
	Imports         []Import      `xml:"import"`
	Elements        []Element     `xml:"element"`
	Attributes      []Attribute   `xml:"attribute"`
	ComplexTypes    []ComplexType `xml:"complexType"`
}

func Parse(xsdPath string) (*Schema, error) {
	f, err := os.Open(xsdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var schema Schema
	d := xml.NewDecoder(f)

	if err := d.Decode(&schema); err != nil {
		return nil, fmt.Errorf("Error decoding XSD: %s", err)
	}

	dir := filepath.Dir(xsdPath)
	for idx, _ := range schema.Imports {
		if err := schema.Imports[idx].Load(dir); err != nil {
			return nil, err
		}
	}
	schema.compile()

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
	return innerSchema.GetElement(ref.Name())
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
			return imp.importedSchema
		}
	}
	return nil
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

func (sch *Schema) GoPackageName() string {
	return sch.Xmlns.PrefixByUri(sch.TargetNamespace)
}

func (sch *Schema) GoImportsNeeded() []string {
	return []string{"encoding/xml"}
}

type Import struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema import"`
	Namespace      string   `xml:"namespace,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	importedSchema *Schema  `xml:"-"`
}

func (i *Import) Load(baseDir string) (err error) {
	i.importedSchema, err = Parse(filepath.Join(baseDir, i.SchemaLocation))
	return
}
