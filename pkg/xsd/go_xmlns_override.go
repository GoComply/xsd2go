package xsd

import (
	"fmt"
	"strings"
)

type xmlnsOverrides map[string]string

func ParseXmlnsOverrides(overrides []string) (xmlnsOverrides, error) {
	ret := xmlnsOverrides{}
	for _, override := range overrides {
		fields := strings.FieldsFunc(override, func(r rune) bool {
			return r == '='
		})
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid xmlns override: '%s' expecting exactly one '=' in the string", override)
		}
		ret[fields[0]] = fields[1]
	}
	return ret, nil
}

func (xo xmlnsOverrides) override(xmlns string) string {
	return xo[xmlns]
}
