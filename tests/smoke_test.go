package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gocomply/xsd2go/pkg/xsd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanity(t *testing.T) {
	xsdFiles, err := filepath.Glob("xsd-examples/valid/*.xsd")
	assert.Nil(t, err)
	assert.NotEmpty(t, xsdFiles)

	for _, xsdPath := range xsdFiles {
		assertConvertsFine(t, xsdPath)
	}
}

func TestSequenceWithinChoice(t *testing.T) {
	xsdPath := "xsd-examples/valid/complex.xsd"
	actual := assertConvertsFine(t, xsdPath)
	expected, err := ioutil.ReadFile(xsdPath + ".out")
	require.NoError(t, err)
	assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), string(actual))
}

func TestRestriction(t *testing.T) {
	xsdPath := "xsd-examples/valid/restriction.xsd"
	actual := assertConvertsFine(t, xsdPath)
	expected, err := ioutil.ReadFile(xsdPath + ".out")
	require.NoError(t, err)
	assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), string(actual))
}

func assertConvertsFine(t *testing.T, xsdPath string) []byte {
	dname, err := ioutil.TempDir("", "xsd2go_tests_")
	assert.Nil(t, err)
	defer os.RemoveAll(dname)

	outputDir := dname

	goModule := "user.com/private"

	err = xsd2go.Convert(xsdPath, goModule, outputDir)
	require.NoError(t, err)

	result, err := ioutil.ReadFile(filepath.Join(outputDir, "simple_schema", "models.go"))
	require.NoError(t, err)

	return result

}
