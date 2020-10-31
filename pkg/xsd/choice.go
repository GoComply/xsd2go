package xsd

import (
	"encoding/xml"
)

type Choice struct {
	XMLName   xml.Name  `xml:"http://www.w3.org/2001/XMLSchema choice"`
	MinOccurs string    `xml:"minOccurs,attr"`
	MaxOccurs string    `xml:"maxOccurs,attr"`
	Elements  []Element `xml:"element"`
	Sequence  *Sequence `xml:"sequence"`
	schema    *Schema   `xml:"-"`
}

func (c *Choice) compile(sch *Schema, parentElement *Element) {
	c.schema = sch
	for idx, _ := range c.Elements {
		el := &c.Elements[idx]

		el.compile(sch, parentElement)
		// Propagate array cardinality downwards
		if c.MaxOccurs == "unbounded" {
			el.MaxOccurs = "unbounded"
		}
		if el.MinOccurs == "" {
			el.MinOccurs = "0"
		}
	}
	if c.Sequence != nil {
		el := c.Sequence
		el.compile(sch, parentElement)
		for _, el2 := range el.Elements() {
			if c.MaxOccurs == "unbounded" {
				el2.MaxOccurs = "unbounded"
			}
			if el2.MinOccurs == "" {
				el2.MinOccurs = "0"
			}
			c.Elements = deduplicateElements(append(c.Elements, el2))
		}
	}

}
