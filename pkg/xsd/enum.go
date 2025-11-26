package xsd

import (
	"encoding/xml"
	"strings"

	"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute.
type Enumeration struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Value   string   `xml:"value,attr"`
}

// Public Go Name of this struct item.
func (e *Enumeration) GoName() string {
	final := e.Value

	// Touch ups for special characters.
	replacer := strings.NewReplacer(
		"+", "Plus",
		"-", "Minus",
		" ", "_",
		".", "_",
		":", "_",
		"/", "_",
		"\\", "_",
		",", "_",
	)
	final = replacer.Replace(final)

	// Only camel case if the enum was at least 3 characters or more to prevent collisions when enums like `HK`, `hk`
	// and `Hk` are all present. Without hints, there is no elegant way to handle this special case.
	if len(e.Value) > 2 {
		final = strcase.ToCamel(strings.ToLower(final))
	}

	return final
}

func (*Enumeration) Modifiers() string {
	return "-"
}

func (e *Enumeration) XmlName() string {
	return e.Value
}
