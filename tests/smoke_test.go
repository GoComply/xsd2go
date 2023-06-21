package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moov-io/xsd2go/pkg/xsd2go"
)

var tests = []struct {
	xsdFile         string
	outputDir       string
	outputFile      string
	goPackage       string
	namespacePrefix string
	expectedFiles   []string
}{
	{
		xsdFile:         "complex.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "complex",
		expectedFiles:   []string{"complex.go.out"},
	},
	{
		xsdFile:         "cpe-naming_2.3.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "cpe-naming_2.3",
		expectedFiles:   []string{"cpe-naming_2.3.go.out"},
	},
	{
		xsdFile:         "restriction.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "restriction",
		expectedFiles:   []string{"restriction.go.out"},
	},
	{
		xsdFile:         "simple.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "simple",
		expectedFiles:   []string{"simple.go.out"},
	},
	{
		xsdFile:         "simple-8859-1.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "simple-8859-1",
		expectedFiles:   []string{"simple-8859-1.go.out"},
	},
	{
		xsdFile:         "swid-2015-extensions-1.0.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "swid-2015-extensions-1.0",
		expectedFiles:   []string{"swid-2015-extensions-1.0.go.out"},
	},
	{
		xsdFile:         "xmldsig-core-schema.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "xmldsig-core-schema",
		expectedFiles:   []string{"xmldsig-core-schema.go.out"},
	},
	{
		xsdFile:         "incl.xsd",
		outputDir:       "simple_schema",
		outputFile:      "models.go",
		goPackage:       "user.com/private",
		namespacePrefix: "incl",
		expectedFiles:   []string{"incl.go.out"},
	},
}

func TestSanity(t *testing.T) {
	dname, err := os.MkdirTemp("", "xsd2go_tests_")
	assert.Nil(t, err)
	defer os.RemoveAll(dname)

	xsdPath := "xsd-examples/xsd/"
	expectedPath := "xsd-examples/assertions"

	for indx := range tests {
		t.Run(tests[indx].xsdFile, func(t *testing.T) {
			outputDir := path.Join(dname, tests[indx].outputDir)
			xsdFile := path.Join(xsdPath, tests[indx].xsdFile)

			err = xsd2go.Convert(
				xsdFile,
				outputDir,
				tests[indx].outputFile,
				tests[indx].goPackage,
				tests[indx].namespacePrefix,
				"rtp",
			)
			require.NoError(t, err)

			golangFiles, err := filepath.Glob(outputDir + "/*")
			require.NoError(t, err)
			assert.Equal(t, len(tests[indx].expectedFiles), len(golangFiles), "Expected to find %v generated files in %s but found %v", len(tests[indx].expectedFiles), outputDir, len(golangFiles))

			for indx2 := range tests[indx].expectedFiles {
				if indx2 < len(golangFiles) {
					actual, err := os.ReadFile(golangFiles[indx2])
					require.NoError(t, err)

					expected, err := os.ReadFile(path.Join(expectedPath, tests[indx].expectedFiles[indx2]))
					require.NoError(t, err)

					t.Logf("Comparing %s to %s", golangFiles[indx2], tests[indx].expectedFiles[indx2])
					assert.Equal(t, string(expected), string(actual))
				}
			}
		})
	}
}
