package xsd

import (
	"fmt"
	"strings"
)

type typeOverrides map[string]map[string]string

func parseTypeOverrides(overrides []string) (typeOverrides, error) {
	result := make(typeOverrides)
	for _, override := range overrides {
		// Split the input string on the `=` to get the schema/type and the type to override with.
		parts := strings.SplitN(override, "=", 2)
		if len(parts) != 2 {
			return typeOverrides{}, fmt.Errorf(`invalid override format: expected at least 1 '=' in "%s"`, override)
		}

		// Split the schema/type to get the schema and the type to be overridden separately.
		index := strings.LastIndexByte(parts[0], ':')
		schemaUri := parts[0][:index]
		typeName := parts[0][index+1:]

		// Add to the resulting map.
		overrideMap, ok := result[schemaUri]
		if !ok {
			overrideMap = make(map[string]string)
		}
		overrideMap[typeName] = parts[1]
		result[schemaUri] = overrideMap
	}
	return result, nil
}

func (t typeOverrides) overrideType(schema string, typeName string) (string, bool) {
	overrides, ok := t[schema]
	if !ok {
		return "", false
	}
	return overrides[typeName], true
}
