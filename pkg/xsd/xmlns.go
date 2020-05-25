package xsd

import (
	"encoding/xml"
)

type xmlns struct {
	Prefix string
	Uri    string
}

type Xmlns []xmlns

func parseXmlns(el xml.StartElement) (result Xmlns) {
	for _, attr := range el.Attr {
		if attr.Name.Space == "xmlns" {
			result = append(result, xmlns{
				Prefix: attr.Name.Local,
				Uri:    attr.Value,
			})
		}
	}
	return
}

func (declarations Xmlns) PrefixByUri(uri string) string {
	for _, declaration := range declarations {
		if declaration.Uri == uri {
			return declaration.Prefix
		}
	}
	return ""
}
