package xsd

import (
	"encoding/xml"
)

type Sequence struct {
	XMLName  xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Elements []Element `xml:"element"`
	Choices  []Choice  `xml:"choice"`
}

func (s *Sequence) compile(sch *Schema) {
	for idx, _ := range s.Elements {
		el := &s.Elements[idx]
		el.compile(sch)
	}
	for idx, _ := range s.Choices {
		c := &s.Choices[idx]
		c.compile(sch)
	}
}
