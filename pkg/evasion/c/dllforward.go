package c

import (
	"embed"
	"fmt"
	"os"

	"github.com/guervild/uru/pkg/models"
	peparser "github.com/saferwall/pe"
)

var outFilePath = "data/templates/c/evasions/dllforward/example.def"

type CDllForwardEvasion struct {
	Name         string
	Debug        bool
	File         string
	ExpectedPath string
	Description  string
}

func NewCDllForwardEvasion() models.ObjectModel {
	return &CDllForwardEvasion{
		Name:  "dllforward",
		Debug: false,
		Description: `Forward exported functions to correct dll. Must be full path
  Argument(s):
    File: Local path to file you want to impersonate
    ExpectedPath: specify the files typical place on machine (including file name), if is not the same as above`,
	}
}

func (e *CDllForwardEvasion) GetImports() []string {
	return []string{}
}

func (e *CDllForwardEvasion) RenderInstanciationCode(data embed.FS) (string, error) {
	// format export strings
	exportLines, err := getNameAndOrdinal(e.File, e.ExpectedPath)

	// handle to file to get fowards
	f, err := os.Create(outFilePath)
	if err != nil {
		return "", fmt.Errorf("Failed to open export definition file")
	}

	// write export lines to outfile
	f.WriteString("EXPORTS\n")
	for _, line := range exportLines {
		f.WriteString(line)
	}

	return "", nil
}

func (e *CDllForwardEvasion) RenderFunctionCode(data embed.FS) (string, error) {
	return "", nil
}

func getNameAndOrdinal(file string, expectedPath string) ([]string, error) {
	// out string array
	var outArray []string

	// parse inpute pe
	pe, err := peparser.New(file, &peparser.Options{})
	if err != nil {
		return nil, fmt.Errorf("Error while opening file: %s, reason: %w", file, err)
	}

	err = pe.Parse()
	if err != nil {
		return nil, fmt.Errorf("Error while parsing file: %s, reason: %w", file, err)
	}

	// check if expected path == file path
	if expectedPath == "" {
		expectedPath = file
	}

	// format each export
	for _, exportedFunction := range pe.Export.Functions {
		outArray = append(outArray, fmt.Sprintf("\t%s=\"%s.%s\" @%d\n", exportedFunction.Name,
			expectedPath, exportedFunction.Name, exportedFunction.Ordinal))
	}

	return outArray, nil
}
