package xsd

import (
	"encoding/xml"
)

type Extension struct {
	XMLName          xml.Name         `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base             reference        `xml:"base,attr"`
	AttributesDirect []Attribute      `xml:"attribute"`
	AttributeGroups  []AttributeGroup `xml:"attributeGroup"`
	Sequence         *Sequence        `xml:"sequence"`
	typ              Type
}

func (ext *Extension) Attributes() []Attribute {
	attrs := ext.AttributesDirect
	if ext.typ != nil {
		attrs = append(attrs, ext.typ.Attributes()...)
		attrs = deduplicateAttributes(attrs)
	}
	for idx := range ext.AttributeGroups {
		attrGroup := ext.AttributeGroups[idx]
		attrs = append(attrs, attrGroup.Attributes()...)
		attrs = deduplicateAttributes(attrs)
	}
	return attrs
}

func (ext *Extension) Elements() []Element {
	elements := []Element{}
	if ext.Sequence != nil {
		elements = append(elements, ext.Sequence.Elements()...)
	}
	if ext.typ != nil {
		elements = append(elements, ext.typ.Elements()...)
		elements = deduplicateElements(elements)
	}

	attrs := ext.Attributes()
	goNames := make(map[string]struct{}, len(elements)+len(attrs))
	for _, attr := range attrs {
		goNames[attr.GoName()] = struct{}{}
	}

	final := []Element{}
	for _, element := range elements {
		if _, found := goNames[element.GoFieldName()]; found {
			element.FieldOverride = true
		}
		goNames[element.GoFieldName()] = struct{}{}
		final = append(final, element)
	}
	return final
}

func deduplicateAttributes(attributes []Attribute) []Attribute {
	seen := make(map[string]struct{}, len(attributes))
	j := 0
	for _, attribute := range attributes {
		if _, ok := seen[attribute.GoName()]; ok {
			continue
		}
		seen[attribute.GoName()] = struct{}{}
		attributes[j] = attribute
		j++
	}
	return attributes[:j]
}

func deduplicateElements(elements []Element) []Element {
	seen := make(map[string]struct{}, len(elements))
	j := 0
	for _, element := range elements {
		if _, ok := seen[element.GoName()]; ok {
			continue
		}
		seen[element.GoName()] = struct{}{}
		elements[j] = element
		j++
	}
	return elements[:j]
}

func (ext *Extension) ContainsText() bool {
	return ext.Base == "xsd:string" || (ext.typ != nil && ext.typ.ContainsText())
}

func (ext *Extension) compile(sch *Schema, parentElement *Element) {
	if ext.Sequence != nil {
		ext.Sequence.compile(sch, parentElement)
	}
	if ext.Base == "" {
		panic("Not implemented: xsd:extension/@base empty, cannot extend unknown type")
	}

	ext.typ = sch.findReferencedType(ext.Base)
	if ext.typ == nil {
		panic("Cannot build xsd:extension: unknown type: " + string(ext.Base))
	}
	ext.typ.compile(sch, parentElement)

	for idx := range ext.AttributeGroups {
		attrGroup := &ext.AttributeGroups[idx]
		attrGroup.compile(sch, parentElement)

	}

	// Handle improbable name clash. Consider XSD defining two attributes on the element:
	// "id" and "Id", this would create name clash given the camelization we do.
	goNames := map[string]uint{}
	for idx := range ext.Attributes() {
		attribute := &ext.Attributes()[idx]
		attribute.compile(sch)

		count := goNames[attribute.GoName()]
		count += 1
		goNames[attribute.GoName()] = count
		attribute.DuplicateCount = count
		// Second GoName may be different depending on the DuplicateCount
		goNames[attribute.GoName()] = count
	}
}
