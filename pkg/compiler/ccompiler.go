// modified from https:// github.com/BishopFox/sliver/blob/5bcfa4c249341e9c9032abcaaf1d4cf459e20059/server/gogo/go.go

package compiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/logger"
)

type CConfig struct {
	ProjectDir     string
	TargetOs       string
	TargetCompiler string
	OutDir         string
	CompileFlags   []string
	ExportDefPath  string
	ArtifactList   []string
	Env            []string
}

func (c *CConfig) GetExportNames(extension string) string {
	return ""
}

func (c *CConfig) GetDebugImports() []string {
	return []string{}
}

func NewEmptyCConfig() *CConfig {
	config := &CConfig{}
	return config
}

func (c *CConfig) PrepareBuild(buildData BuildData) error {
	targetCompiler := "x86_64-w64-mingw32-gcc"
	var compileFlags []string

	if common.ContainsStringInSliceIgnoreCase(buildData.Imports, "iostream") {
		// iostream is a c++ library
		targetCompiler = "x86_64-w64-mingw32-g++"

		// this is a lazy was to solve several problems
		// conversion from const char * to unsigned char * (in c++ string literals are const char arrays and
		// we nee this flag to convert ot unsigned char array, there are other ways to do, but this is easy)
		compileFlags = append(compileFlags, "-fpermissive")
	}

	path, err := IsTargetCompilerInstalled(targetCompiler)

	if err != nil {
		return fmt.Errorf("target compiler not found %w", err)
	}

	logger.Logger.Debug().Str("target_compiler", path).Msg("Path to the target compiler")

	// -static needed because we cant assume mingw will be on the target system
	compileFlags = append(compileFlags, "--static")

	// -static-libgcc and -static-libstdc++ needed to link the C and C++ standard libraries statically and
	// remove the need to carry around any separate copies of those.
	compileFlags = append(compileFlags, "-static-libstdc++", "-static-libgcc")

	if buildData.BuildMode == "c-shared" {
		compileFlags = append(compileFlags, "-shared")
	}

	c.ProjectDir = buildData.DirPath
	c.TargetOs = buildData.TargetOs
	c.TargetCompiler = targetCompiler
	c.CompileFlags = compileFlags
	c.ArtifactList = buildData.ArtifactList

	return nil
}

func (c *CConfig) Build(payload, dest string) ([]byte, error) {
	// remove abs path from dest var
	c.CompileFlags = append(c.CompileFlags, fmt.Sprintf("-o%s", filepath.Base(dest)))

	// format arguments to be used by compiler
	c.CompileFlags = append(c.CompileFlags, filepath.Base(payload))

	if common.ContainsStringInSliceIgnoreCase(c.ArtifactList, "dllforward") {
		c.CompileFlags = append(c.CompileFlags, []string{"../../data/templates/c/evasions/dllforward/example.def"}...)
	}
	logger.Logger.Debug().Str("project_dir", c.ProjectDir).Msg("Project dir")
	logger.Logger.Debug().Str("compile_args", strings.Join(c.CompileFlags, " ")).Msg("Defining compile arguments")

	// create command (dont run it yet)
	cmd := exec.Command(c.TargetCompiler, c.CompileFlags...)

	cmd.Dir = c.ProjectDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run command, gather output
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error %w: %s", err, string(stderr.String()))
	}

	return stdout.Bytes(), err
}

func (c *CConfig) IsTypeSupported(t string) (string, string, error) {
	switch strings.ToLower(t) {
	case "exe":
		return "exe", "", nil
	case "dll":
		return "dll", "c-shared", nil
	default:
		return "", "", fmt.Errorf("unsupported executable type")
	}
}