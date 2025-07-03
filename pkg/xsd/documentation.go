package xsd

import (
	"encoding/xml"
	"regexp"
	"strings"
)

type Documentation struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema documentation"`
	Source  string   `xml:"source,attr"`
	Lang    string   `xml:"xml:lang,attr"`
	Content []byte   `xml:",innerxml"`
}

func (d Documentation) GetContent() string {
	docs := string(d.Content)
	docs = regexp.MustCompile(`\s`).ReplaceAllString(docs, " ")
	docs = regexp.MustCompile(" +").ReplaceAllString(docs, " ")
	docs = strings.Trim(docs, " ")
	return docs
}
