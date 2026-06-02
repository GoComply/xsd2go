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
	"golang.org/x/tools/txtar"
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

	outputDir := t.TempDir()
	goModule := "user.com/private"

	err := xsd2go.Convert(xsdPath, goModule, outputDir, nil)
	require.NoError(t, err)

	generatedFilePath, err := locateGeneratedFile(outputDir)
	require.NoError(t, err)
	result, err := os.ReadFile(generatedFilePath)
	require.NoError(t, err)

	out, err := exec.CommandContext(t.Context(), "go", "build", generatedFilePath).CombinedOutput()
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

func extractTxtar(t *testing.T, name string) (dir string) {
	t.Helper()

	data, err := os.ReadFile(name)
	require.NoError(t, err)

	ar := txtar.Parse(data)

	dir = os.TempDir()

	for _, f := range ar.Files {
		path := filepath.Join(dir, f.Name)

		require.NoError(t,
			os.MkdirAll(filepath.Dir(path), 0755))

		require.NoError(t,
			os.WriteFile(path, f.Data, 0600))
	}

	return
}

func TestCircularImport(t *testing.T) {
	workdir, err := os.Getwd()
	require.NoError(t, err)

	xsdFiles, err := filepath.Glob("xsd-examples/modules/*.txtar")
	require.NoError(t, err)
	assert.NotEmpty(t, xsdFiles)

	for _, xsdPath := range xsdFiles {
		dir := extractTxtar(t, filepath.Join(workdir, xsdPath))
		t.Chdir(dir)

		data, err := os.ReadFile(filepath.Join(workdir, xsdPath+".out"))
		require.NoError(t, err)
		expectar := txtar.Parse(data)

		err = xsd2go.Convert(
			"a.xsd",
			"example.com/test",
			"out",
			nil,
		)
		require.NoError(t, err)

		for _, expect := range expectar.Files {
			got, err := os.ReadFile(filepath.Join(dir, "out", expect.Name))
			require.NoError(t, err)

			assert.Equal(t, strings.ReplaceAll(string(expect.Data), "\r\n", "\n"), string(got))
		}
	}
}
