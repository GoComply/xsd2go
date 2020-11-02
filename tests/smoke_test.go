package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
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
		actual := assertConvertsFine(t, xsdPath)

		expected, err := ioutil.ReadFile(xsdPath + ".out")
		require.NoError(t, err)
		assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), string(actual))
	}
}

func assertConvertsFine(t *testing.T, xsdPath string) []byte {
	dname, err := ioutil.TempDir("", "xsd2go_tests_")
	assert.Nil(t, err)
	defer os.RemoveAll(dname)

	outputDir := dname

	goModule := "user.com/private"

	err = xsd2go.Convert(xsdPath, goModule, outputDir)
	require.NoError(t, err)

	generatedFile := filepath.Join(outputDir, "simple_schema", "models.go")
	result, err := ioutil.ReadFile(generatedFile)
	require.NoError(t, err)

	out, err := exec.Command("go", "build", generatedFile).Output()
	require.NoError(t, err)
	assert.Equal(t, string(out), "")

	return result

}
