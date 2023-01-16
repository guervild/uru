package compiler

import (
	"embed"
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
func GetEmptyCompiler(lang string) compiler {

	switch lang {
	case "go":
		return NewEmptyGoConfig()
	case "c":
		return NewEmptyCConfig()
	default:
		return nil
	}
	return nil
}

func GetSupportedLangs(lang string) bool {
	return contains(supportedLangs, lang)
}
