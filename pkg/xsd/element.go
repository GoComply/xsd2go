package xsd

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/iancoleman/strcase"
)

// Element defines single XML element.
type Element struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name            string   `xml:"name,attr"`
	nameOverride    string
	XmlNameOverride string      `xml:"-"`
	FieldOverride   bool        `xml:"-"`
	Type            reference   `xml:"type,attr"`
	Ref             reference   `xml:"ref,attr"`
	MinOccurs       string      `xml:"minOccurs,attr"`
	MaxOccurs       string      `xml:"maxOccurs,attr"`
	Annotation      *Annotation `xml:"annotation"`
	refElm          *Element
	ComplexType     *ComplexType `xml:"complexType"`
	SimpleType      *SimpleType  `xml:"simpleType"`
	schema          *Schema
	typ             Type
	Parent          *Element `xml:"-"`
}

func (e *Element) Attributes() []Attribute {
	if e.typ != nil {
		return injectSchemaIntoAttributes(e.schema, e.typ.Attributes())
	}
	return []Attribute{}
}

func (e *Element) ContainsDocumentation() bool {
	return e.Documentation() != ""
}

func (e *Element) Documentation() string {
	if e.Annotation == nil {
		return ""
	}
	if len(e.Annotation.Documentations) == 0 {
		return ""
	}
	return e.Annotation.Documentations[0].GetContent()
}

func (e *Element) Elements() []Element {
	if e.typ != nil {
		return e.typ.Elements()
	}
	return []Element{}
}

func (e *Element) GoFieldName() string {
	name := e.Name
	if name == "" {
		return e.refElm.GoName()
	}
	if e.FieldOverride {
		name += "Elm"
	}
	return strcase.ToCamel(name)
}

func (e *Element) GoName() string {
	if e.nameOverride != "" {
		return strcase.ToCamel(e.nameOverride)
	}
	return e.GoFieldName()
}

func (e *Element) GoMemLayout() string {
	if e.isArray() {
		return "[]"
	}
	if (e.MaxOccurs == "1" || e.MaxOccurs == "") && e.MinOccurs == "0" && e.GoTypeName() != "string" {
		return "*"
	}
	return ""
}

func (e *Element) GoTypeName() string {
	if e.Type != "" {
		return e.typ.GoName()
	} else if e.Ref != "" {
		return e.refElm.GoTypeName()
	} else if e.isPlainString() {
		return "string"
	}
	return e.GoName()
}

func (e *Element) GoForeignModule() string {
	if e.isPlainString() && e.refElm == nil && e.typ == nil {
		return ""
	}

	foreignSchema := (*Schema)(nil)
	if e.refElm != nil {
		foreignSchema = e.refElm.schema
	} else if e.typ != nil {
		foreignSchema = e.typ.Schema()
	}

	if foreignSchema != nil && foreignSchema != e.schema &&
		foreignSchema.TargetNamespace != e.schema.TargetNamespace {
		return foreignSchema.GoPackageName() + "."
	}
	return ""
}

func (e *Element) Modifiers() string {
	res := ""
	if e.optional() {
		res += ",omitempty"
	}
	return res
}

func (e *Element) optional() bool {
	return e.MinOccurs == "0"
}

func (e *Element) XmlName() string {
	if e.XmlNameOverride != "" {
		return e.XmlNameOverride
	}
	name := e.Name
	if name == "" {
		return e.refElm.XmlName()
	}
	return name
}

func (e *Element) ContainsText() bool {
	return e.typ != nil && e.typ.ContainsText()
}

func (e *Element) isPlainString() bool {
	return e.SimpleType != nil || (e.Type == "" && e.Ref == "" && e.ComplexType == nil) || (e.typ != nil && e.typ.GoTypeName() == "string")
}

func (e *Element) isArray() bool {
	if e.MaxOccurs == "unbounded" {
		return true
	}
	occurs, err := strconv.Atoi(e.MaxOccurs)
	return err == nil && occurs > 1
}

func (e *Element) compile(s *Schema, parentElement *Element) {
	e.schema = s
	e.Parent = parentElement
	if e.ComplexType != nil {
		e.typ = e.ComplexType
		if e.SimpleType != nil {
			panic("Not implemented: xsd:element " + e.Name + " defines ./xsd:simpleType and ./xsd:complexType together")
		} else if e.Type != "" {
			panic("Not implemented: xsd:element " + e.Name + " defines ./@type= and ./xsd:complexType together")
		}
		e.typ.compile(s, e)
	} else if e.SimpleType != nil {
		e.typ = e.SimpleType
		if e.Type != "" {
			panic("Not implemented: xsd:element " + e.Name + " defines ./@type= and ./xsd:simpleType together")
		}
		e.typ.compile(s, e)
	} else if e.Type != "" {
		e.typ = e.schema.findReferencedType(e.Type)
		if e.typ == nil {
			panic("Cannot resolve type reference: " + string(e.Type))
		}
	} else if e.Ref != "" {
		e.refElm = e.schema.findReferencedElement(e.Ref)
		if e.refElm == nil {
			panic("Cannot resolve element reference: " + e.Ref)
		}
	}

	if e.Ref == "" && e.Type == "" && !e.isPlainString() {
		e.schema.registerInlinedElement(e, parentElement)
	}
}

func (e *Element) prefixNameWithParent(parentElement *Element) {
	// In case there are inlined xsd:elements within another xsd:elements, it may happen that two top-level xsd:elements
	// define child xsd:element of a same name. In such case, we need to override children name to avoid name clashes.
	if parentElement != nil {
		e.nameOverride = fmt.Sprintf("%s-%s", parentElement.GoName(), e.GoName())
	}
}
