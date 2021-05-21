package xsd

import (
	"encoding/xml"
	"html"
	"strings"
)

// Attribute defines single XML attribute
type Annotation struct {
	XMLName           xml.Name        `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	DocumentationList []Documentation `xml:"documentation"`
}

func (a *Annotation) GetName() string {
	if a != nil {
		for _, d := range a.DocumentationList {
			if strings.EqualFold("Name", d.Source) {
				return d.InnerXml
			}
		}
	}
	return ""
}

func (a *Annotation) GoComments() []string {
	if a != nil {
		for _, d := range a.DocumentationList {
			if strings.EqualFold("Definition", d.Source) {
				def := html.UnescapeString(d.InnerXml)
				def = strings.ReplaceAll(def, "\r", "")
				return strings.Split(def, "\n")
			}
		}
	}
	return nil
}
