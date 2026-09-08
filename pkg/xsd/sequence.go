package xsd

import (
	"encoding/xml"
)

type Sequence struct {
	XMLName     xml.Name   `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	MinOccurs   string     `xml:"minOccurs,attr"`
	MaxOccurs   string     `xml:"maxOccurs,attr"`
	ElementList []Element  `xml:"element"`
	Choices     []Choice   `xml:"choice"`
	Sequences   []Sequence `xml:"sequence"`
	allElements []Element
}

func (s *Sequence) Elements() []Element {
	return s.allElements
}

func (s *Sequence) compile(sch *Schema, parentElement *Element) {
	for idx := range s.ElementList {
		el := &s.ElementList[idx]
		el.compile(sch, parentElement)
		if s.MaxOccurs == "unbounded" {
			el.MaxOccurs = "unbounded"
		} else if s.MaxOccurs != "" && el.MaxOccurs == "" {
			el.MaxOccurs = s.MaxOccurs
		}
		if s.MinOccurs == "0" && el.MinOccurs == "" {
			el.MinOccurs = "0"
		}
	}
	s.allElements = s.ElementList

	for idx := range s.Choices {
		c := &s.Choices[idx]
		c.compile(sch, parentElement)

		for _, nestedEl := range c.Elements() {
			if s.MaxOccurs == "unbounded" {
				nestedEl.MaxOccurs = "unbounded"
			} else if s.MaxOccurs != "" && nestedEl.MaxOccurs == "" {
				nestedEl.MaxOccurs = s.MaxOccurs
			}
			if s.MinOccurs == "0" && nestedEl.MinOccurs == "" {
				nestedEl.MinOccurs = "0"
			}
		}

		s.allElements = append(s.allElements, c.Elements()...)
	}

	for idx := range s.Sequences {
		nestedSeq := &s.Sequences[idx]
		nestedSeq.compile(sch, parentElement)

		// Propagate cardinality from parent sequence to nested sequence's elements
		for _, nestedEl := range nestedSeq.allElements {
			if s.MaxOccurs == "unbounded" {
				nestedEl.MaxOccurs = "unbounded"
			} else if s.MaxOccurs != "" && nestedEl.MaxOccurs == "" {
				nestedEl.MaxOccurs = s.MaxOccurs
			}
			if s.MinOccurs == "0" && nestedEl.MinOccurs == "" {
				nestedEl.MinOccurs = "0"
			}
		}

		s.allElements = append(s.allElements, nestedSeq.allElements...)
	}
}

type SequenceAll struct {
	XMLName     xml.Name   `xml:"http://www.w3.org/2001/XMLSchema all"`
	MinOccurs   string     `xml:"minOccurs,attr"`
	MaxOccurs   string     `xml:"maxOccurs,attr"`
	ElementList []Element  `xml:"element"`
	Choices     []Choice   `xml:"choice"`
	Sequences   []Sequence `xml:"sequence"`
	allElements []Element
}

func (s *SequenceAll) Elements() []Element {
	return s.allElements
}

func (s *SequenceAll) compile(sch *Schema, parentElement *Element) {
	for idx := range s.ElementList {
		el := &s.ElementList[idx]
		el.compile(sch, parentElement)
		if s.MaxOccurs == "unbounded" {
			el.MaxOccurs = "unbounded"
		} else if s.MaxOccurs != "" && el.MaxOccurs == "" {
			el.MaxOccurs = s.MaxOccurs
		}
		if s.MinOccurs == "0" && el.MinOccurs == "" {
			el.MinOccurs = "0"
		}
	}
	s.allElements = s.ElementList

	for idx := range s.Choices {
		c := &s.Choices[idx]
		c.compile(sch, parentElement)

		for _, nestedEl := range c.Elements() {
			if s.MaxOccurs == "unbounded" {
				nestedEl.MaxOccurs = "unbounded"
			} else if s.MaxOccurs != "" && nestedEl.MaxOccurs == "" {
				nestedEl.MaxOccurs = s.MaxOccurs
			}
			if s.MinOccurs == "0" && nestedEl.MinOccurs == "" {
				nestedEl.MinOccurs = "0"
			}
		}

		s.allElements = append(s.allElements, c.Elements()...)
	}

	for idx := range s.Sequences {
		nestedSeq := &s.Sequences[idx]
		nestedSeq.compile(sch, parentElement)

		// Propagate cardinality from parent sequence to nested sequence's elements
		for _, nestedEl := range nestedSeq.allElements {
			if s.MaxOccurs == "unbounded" {
				nestedEl.MaxOccurs = "unbounded"
			} else if s.MaxOccurs != "" && nestedEl.MaxOccurs == "" {
				nestedEl.MaxOccurs = s.MaxOccurs
			}
			if s.MinOccurs == "0" && nestedEl.MinOccurs == "" {
				nestedEl.MinOccurs = "0"
			}
		}

		s.allElements = append(s.allElements, nestedSeq.allElements...)
	}
}
