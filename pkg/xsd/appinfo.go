package xsd

import "encoding/xml"

type AppInfo struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema appinfo"`
}
