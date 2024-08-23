package xsd

import (
	"encoding/xml"
)

type Annotation struct {
	XMLName        xml.Name        `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	ID             string          `xml:"id,attr"`
	AppInfos       []AppInfo       `xml:"appinfo"`
	Documentations []Documentation `xml:"documentation"`
}
