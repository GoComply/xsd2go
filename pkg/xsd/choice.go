package xsd

import (
	"encoding/xml"
)

type Choice struct {
	XMLName   xml.Name  `xml:"http://www.w3.org/2001/XMLSchema choice"`
	MinOccurs string    `xml:"minOccurs,attr"`
	MaxOccurs string    `xml:"maxOccurs,attr"`
	Elements  []Element `xml:"element"`
	schema    *Schema   `xml:"-"`
}

func (c *Choice) compile(sch *Schema) {
	c.schema = sch
	for idx, _ := range c.Elements {
		el := &c.Elements[idx]
		el.compile(sch)
	}
}
