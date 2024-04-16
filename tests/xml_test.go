package tests

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gocomply/xsd2go/tests/test-schemas/complex"
)

func TestComplex(t *testing.T) {
	var v complex.Myelement
	assertMarshal(&v, "xml-examples/complex.xsd.1.xml", t)
	var z complex.Myelement
	assertMarshal(&z, "xml-examples/complex.xsd.2.xml", t)
	var x complex.Myelement
	assertMarshal(&x, "xml-examples/complex.xsd.3.xml", t)
}

func assertMarshal(v any, xmlPath string, t *testing.T) {
	//unmarshal xml file into v
	in, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Fatalf("Failure opening test file: %v", err)
	}
	err = xml.Unmarshal(in, v)
	if err != nil {
		t.Fatalf("Failure parsing test file: %v", err)
	}
	//marshal v into buff and compare
	out, err := xml.MarshalIndent(v, "", "  ")
	out = append(out, '\n')
	if err != nil {
		t.Fatalf("Failure marshalling output: %v", err)
	}
	expected, err := os.ReadFile(xmlPath)
	if err != nil {
		t.Fatalf("Failure reading result file: %v", err)
	}
	assert.Equal(t, string(expected), string(out))
}
