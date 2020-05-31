package xsd

import (
	"encoding/xml"
)

type Sequence struct {
	XMLName     xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	ElementList []Element `xml:"element"`
	Choices     []Choice  `xml:"choice"`
}

func (s *Sequence) Elements() []Element {
	return s.ElementList
}

func (s *Sequence) compile(sch *Schema) {
	for idx, _ := range s.ElementList {
		el := &s.ElementList[idx]
		el.compile(sch)
	}
	for idx, _ := range s.Choices {
		c := &s.Choices[idx]
		c.compile(sch)
	}
}
