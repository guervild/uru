package compiler

import (
	"embed"
	"fmt"
	"strings"
	"os/exec"

	"github.com/guervild/uru/pkg/common"
)

type compiler interface {
	Build(payload, dest string) ([]byte, error)
	IsTypeSupported(t string) (string, string, error)
	PrepareBuild(buildData BuildData) error
	GetDebugImports() []string
	GetExportNames(extension string) string
}

type BuildData struct {
	TargetOs     string
	Arch         string
	DirPath      string
	BuildMode    string
	Keep         bool
	Trimpath     bool
	Obfuscation  bool
	Imports      []string
	ArtifactList []string
	DataTemplate embed.FS
	FileProps    string
}

var supportedLangs = []string{
	"go",
	"c",
}

// GetEmptyCompiler retrieve correct compiler based on language
func GetEmptyCompiler(lang string) (compiler, error) {

	switch strings.ToLower(lang) {
	case "go":
		return NewEmptyGoConfig(), nil
	case "c":
		return NewEmptyCConfig(), nil
	default:
	}

	return nil, fmt.Errorf("Error ... Language is not supported yet")
}

func GetSupportedLangs(lang string) bool {
	return common.ContainsStringInSliceIgnoreCase(supportedLangs, strings.ToLower(lang))
}

func IsTargetCompilerInstalled(target string) (string, error) {
	return exec.LookPath(target)
}


func GetProperArch(arch string, lang string) (string, error) {

	switch lang {
	case "go":
		if arch == "x64" {
			return "amd64", nil
		} else if arch == "x86" {
			return "386", nil
		}
	case "c":
		if arch == "x64" {
			return "x64", nil
		}
	}

	return "", fmt.Errorf("Arch value must either x86 either x64.")
}

func GetCoreFile(lang string) (string, error) {
	switch lang {
	case "go":
		return "templates/go/core.go.tmpl", nil
	case "c":
		return "templates/c/core.c.tmpl", nil
	}
	return "", fmt.Errorf("golang and c are the only supported languages")
}