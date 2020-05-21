package xsd

import (
	"strings"
)

// Internal XSD reference. Examples: "xml:lang", "cpe2:platform-specification"
type reference string

func (ref reference) NsPrefix() string {
	colonPos := strings.Index(string(ref), ":")
	if colonPos == -1 {
		return ""
	}
	return string(ref)[0:colonPos]
}

func (ref reference) Name() string {
	colonPos := strings.Index(string(ref), ":")
	return string(ref)[colonPos+1:]
}
