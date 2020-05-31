package xsd

import (
	"encoding/xml"
)

type Sequence struct {
	XMLName     xml.Name  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	ElementList []Element `xml:"element"`
	Choices     []Choice  `xml:"choice"`
	allElements []Element `xml:"-"`
}

func (s *Sequence) Elements() []Element {
	return s.allElements
}

func (s *Sequence) compile(sch *Schema) {
	for idx, _ := range s.ElementList {
		el := &s.ElementList[idx]
		el.compile(sch)
	}

	s.allElements = s.ElementList
	for idx, _ := range s.Choices {
		c := &s.Choices[idx]
		c.compile(sch)

		s.allElements = append(s.allElements, c.Elements...)
	}
}
