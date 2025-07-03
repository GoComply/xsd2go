package tests_test

import (
	"fmt"
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
	require.NoError(t, err)
	assert.NotEmpty(t, xsdFiles)

	for _, xsdPath := range xsdFiles {
		actual := assertConvertsFine(t, xsdPath)

		expected, err := os.ReadFile(xsdPath + ".out")
		require.NoError(t, err)
		assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), string(actual))
	}
}

func assertConvertsFine(t *testing.T, xsdPath string) []byte {
	t.Helper()

	dname, err := os.MkdirTemp("", "xsd2go_tests_")
	require.NoError(t, err)
	defer os.RemoveAll(dname)

	outputDir := dname

	goModule := "user.com/private"

	err = xsd2go.Convert(xsdPath, goModule, outputDir, nil)
	require.NoError(t, err)

	generatedFilePath, err := locateGeneratedFile(outputDir)
	require.NoError(t, err)
	result, err := os.ReadFile(generatedFilePath)
	require.NoError(t, err)

	out, err := exec.Command("go", "build", generatedFilePath).CombinedOutput()
	assert.Empty(t, string(out))
	require.NoError(t, err)

	return result
}

func locateGeneratedFile(outputDir string) (string, error) {
	golangFiles, err := filepath.Glob(outputDir + "/*/models.go")
	if err != nil {
		return "", err
	}
	if len(golangFiles) != 1 {
		return "", fmt.Errorf("Expected to find single generated file but found %s", golangFiles)
	}
	return golangFiles[0], nil
}
