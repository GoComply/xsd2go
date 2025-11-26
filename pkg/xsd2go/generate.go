package xsd2go

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gocomply/xsd2go/pkg/template"
	"github.com/gocomply/xsd2go/pkg/xsd"
)

func Convert(inputs []string, goModule, outputDir string, xmlnsOverrides []string, typeOverrides []string) error {
	var allInputs []string
	for _, input := range inputs {
		fnames, err := filepath.Glob(input)
		if err != nil {
			return err
		}
		allInputs = append(allInputs, fnames...)
	}

	if len(allInputs) == 0 {
		fmt.Println("No inputs to convert")
		return nil
	}

	// Sanitize inputs.
	for i, input := range allInputs {
		// Trim any unneeded whitespace characters.
		allInputs[i] = strings.TrimSpace(input)

		// Remove unnecessary characters at the beginning of the path.
		allInputs[i] = filepath.Clean(input)

		// Use `/` characters to prevent problems with building paths later.
		allInputs[i] = filepath.ToSlash(input)
	}

	// Report for sanity.
	fmt.Printf("Processing:\n")
	for _, input := range allInputs {
		fmt.Printf("\t%s\n", input)
	}
	fmt.Printf("\n")

	// Build the package output path and use forward slashes so package paths are generated correctly on all operating
	// systems.
	pkgOutputPath := filepath.Join(goModule, outputDir)
	pkgOutputPath = filepath.ToSlash(pkgOutputPath)

	// Create the workspace.
	ws, err := xsd.NewWorkspace(pkgOutputPath, allInputs, xmlnsOverrides, typeOverrides)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}
		if err := template.GenerateTypes(sch, outputDir); err != nil {
			return err
		}
	}

	return nil
}
