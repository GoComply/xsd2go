package xsd

import (
	"encoding/xml"
)

type Choice struct {
	XMLName     xml.Name   `xml:"http://www.w3.org/2001/XMLSchema choice"`
	MinOccurs   string     `xml:"minOccurs,attr"`
	MaxOccurs   string     `xml:"maxOccurs,attr"`
	ElementList []Element  `xml:"element"`
	Sequences   []Sequence `xml:"sequence"`
	schema      *Schema    `xml:"-"`
	allElements []Element  `xml:"-"`
}

func (c *Choice) compile(sch *Schema, parentElement *Element) {
	c.schema = sch
	for idx := range c.ElementList {
		el := &c.ElementList[idx]

		el.compile(sch, parentElement)
		// Propagate array cardinality downwards
		if c.MaxOccurs == "unbounded" {
			el.MaxOccurs = "unbounded"
		}
		if el.MinOccurs == "" {
			el.MinOccurs = "0"
		}
	}

	c.allElements = c.ElementList
	inheritedElements := []Element{}
	for idx := range c.Sequences {
		el := &c.Sequences[idx]
		el.compile(sch, parentElement)
		for _, el2 := range el.Elements() {
			if c.MaxOccurs == "unbounded" {
				el2.MaxOccurs = "unbounded"
			}
			if el2.MinOccurs == "" {
				el2.MinOccurs = "0"
			}
			inheritedElements = append(inheritedElements, el2)
		}
	}
	// deduplicate elements that represent duplicate within xsd:choice/xsd:sequence structure
	c.allElements = append(c.ElementList, deduplicateElements(inheritedElements)...)

}

func (c *Choice) Elements() []Element {
	return c.allElements
}
