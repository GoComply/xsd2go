package xsd

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// Reference XSD reference. Examples: "xml:lang", "cpe2:platform-specification"
type Reference string

func (ref Reference) NsPrefix() string {
	colonPos := strings.Index(string(ref), ":")
	if colonPos == -1 {
		return ""
	}
	return string(ref)[0:colonPos]
}

func (ref Reference) Name() string {
	colonPos := strings.Index(string(ref), ":")
	return string(ref)[colonPos+1:]
}

func (ref Reference) GoName() string {
	return strcase.ToCamel(ref.NsPrefix()) + strcase.ToCamel(ref.Name())
}
