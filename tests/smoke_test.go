package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gocomply/xsd2go/pkg/xsd2go"
	"github.com/stretchr/testify/assert"
)

func TestSanity(t *testing.T) {
	xsdFiles, err := filepath.Glob("xsd-examples/valid/*.xsd")
	assert.Nil(t, err)
	assert.NotEmpty(t, xsdFiles)

	for _, xsdPath := range xsdFiles {
		assertConvertsFine(t, xsdPath)
	}
}

func assertConvertsFine(t *testing.T, xsdPath string) {
	dname, err := ioutil.TempDir("", "xsd2go_tests_")
	assert.Nil(t, err)
	defer os.RemoveAll(dname)

	outputDir := dname

	goModule := "user.com/private"

	err = xsd2go.Convert(xsdPath, goModule, outputDir)
}
