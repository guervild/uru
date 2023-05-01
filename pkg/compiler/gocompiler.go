// modified from https://github.com/BishopFox/sliver/blob/5bcfa4c249341e9c9032abcaaf1d4cf459e20059/server/gogo/go.go

package compiler

import (
	"bytes"
	"fmt"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/tampering"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/guervild/uru/pkg/logger"
)

type GoConfig struct {
	ProjectDir   string
	GOOS         string
	GOARCH       string
	GOCACHE      string
	Buildmode    string
	Keep         bool
	Trimpath     bool
	Obfuscation  bool
	Path         string
	CompilerPath string

	Imports []string
	Command []string
	Env     []string
}

func (g *GoConfig) GetExportNames(extension string) string {

	lowerExtension := strings.ToLower(extension)

	if lowerExtension == "cpl" {
		return `
		//export CPlApplet
		func CPlApplet() {
			Start()
		}`
	}

	if lowerExtension == "dll" {
		return `
		//export DllRegisterServer
		func DllRegisterServer() {
			Start()
		}
		
		//export DllGetClassObject
		func DllGetClassObject() {
			Start()
		}
		
		//export DllUnregisterServer
		func DllUnregisterServer() {
			Start()
		}`
	}

	if lowerExtension == "xll" {
		return `
		//export xlAutoOpen
		func xlAutoOpen() {
			Start()
		}`
	}

	return ""
}

type GoMod struct {
	RandomName string
}

func NewEmptyGoConfig() *GoConfig {
	config := &GoConfig{}
	return config
}

func (g *GoConfig) GetDebugImports() []string {
	return []string{
		`"io"`,
		`"os"`,
		`"fmt"`,
	}
}

func (g *GoConfig) PrepareBuild(buildData BuildData) error {

	//Copy go.mod
	fileGoMod, err := common.CreatePayloadFile("go", "mod", buildData.DirPath)
	defer fileGoMod.Close()
	if err != nil {
		return err
	}

	// go.mod template formatting
	tGoMod, err := template.ParseFS(buildData.DataTemplate, "templates/go/go.mod.tmpl")
	if err != nil {
		return err
	}
	err = tGoMod.Execute(fileGoMod, &GoMod{RandomName: common.RandomStringOnlyChar(8)})
	if err != nil {
		return err
	}

	//Copy go.sum
	goSum, _ := buildData.DataTemplate.ReadFile("templates/go/go.sum.tmpl")
	if err := os.WriteFile(path.Join(buildData.DirPath, "go.sum"), goSum, 0666); err != nil {
		return err
	}

	//FileProperties
	if buildData.FileProps != "" {
		logger.Logger.Info().Msgf("Try to use file properties: %s", buildData.FileProps)

		name, err := tampering.BuildFromJson(buildData.FileProps, buildData.Arch, buildData.DirPath)

		if err != nil {
			logger.Logger.Info().Msgf("Could not use the file properties: %s", err.Error())
			logger.Logger.Info().Msg("Continue ...")
		} else {
			logger.Logger.Info().Msgf("Successfully used file properties with internal name: %s", name)
		}
	}

	g.ProjectDir = buildData.DirPath
	g.GOOS = buildData.TargetOs
	g.GOARCH = buildData.Arch
	g.Buildmode = buildData.BuildMode
	g.Trimpath = buildData.Trimpath
	g.Obfuscation = buildData.Obfuscation
	g.Keep = buildData.Keep
	g.Imports = buildData.Imports
	g.GOCACHE = getGoCache(buildData.DirPath)

	return nil
}

func (g *GoConfig) Build(payload, dest string) ([]byte, error) {

	g.BuildGoBuildCommand(payload, dest)

	if 0 >= len(g.Command) {
		return nil, fmt.Errorf("Error while setting golang commands ...")
	}

	logger.Logger.Debug().Str("compile_args", strings.Join(g.Command, " ")).Msg("Defining compile arguments")

	cmd := exec.Command(g.CompilerPath, g.Command...)
	cmd.Dir = g.ProjectDir

	cmd.Env = g.Env

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("error %s: %s", err, string(stderr.Bytes()))
	}

	return stdout.Bytes(), err
}

func (g *GoConfig) BuildGoBuildCommand(payload, dest string) {

	// Get Compiler
	var compilerPath string
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		goroot = build.Default.GOROOT
	}

	homeEnv := os.Getenv("HOME")

	if g.Obfuscation {
		compilerPath = path.Join(path.Join(gopath, "bin"), "garble")
	} else {
		compilerPath = path.Join(path.Join(goroot, "bin"), "go")
	}

	g.CompilerPath = compilerPath

	logger.Logger.Debug().Str("compiler_path", g.CompilerPath).Msg("Set the compiler path")

	//Getting the right env variables
	env := []string{"GOPRIVATE=*", fmt.Sprintf("GOCACHE=%s", g.GOCACHE), fmt.Sprintf("GOPATH=%s", gopath), fmt.Sprintf("GOOS=%s", g.GOOS), fmt.Sprintf("GOARCH=%s", g.GOARCH), fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
	
	if homeEnv != "" {
		env = append(env, fmt.Sprintf("HOME=%s", homeEnv))
	}
	
	if g.Buildmode == "c-shared" {
		env = append(env, "CGO_ENABLED=1")
		env = append(env, "CXX=x86_64-w64-mingw32-g++")
		env = append(env, "CC=x86_64-w64-mingw32-gcc")
	}
	// var gogarble []string
	// if g.Obfuscation {
	// 	for _, v := range g.Imports {
	// 		var out string
	// 		out = strings.TrimLeft(strings.TrimRight(v, "\""), "\"")
	// 		if strings.Contains(out, "\"") {
	// 			out2 := strings.Join(strings.Split(out, "\"")[1:], "\"")
	// 			out = out2
	// 		}
	// 		gogarble = append(gogarble, out)
	// 	}
	// 	env = append(env, fmt.Sprintf("GOGARBLE=%s", strings.Join(gogarble, ",")))
	// }

	g.Env = env

	logger.Logger.Debug().Str("compile_env_vars", strings.Join(g.Env, " ")).Msg("Set the environment variables")

	//Setting up go command
	goCommand := []string{"build"}

	ldflags := []string{"-s -w -buildid="}
	goCommand = append(goCommand, "-ldflags")
	goCommand = append(goCommand, ldflags...)

	if g.Trimpath {
		goCommand = append(goCommand, "-trimpath")
	}

	goCommand = append(goCommand, "-a")
	goCommand = append(goCommand, []string{"-o", dest}...)

	if len(g.Buildmode) > 0 {
		goCommand = append(goCommand, fmt.Sprintf("-buildmode=%s", g.Buildmode))
	}

	if g.Obfuscation {
		goCommand = g.GetGarbleArgs(goCommand)
	}

	//Final command
	g.Command = append(g.Command, goCommand...)
}

func (g *GoConfig) GetGarbleArgs(command []string) []string {

	var goCommand = []string{}

	if g.Keep {
		goCommand = append(goCommand, fmt.Sprintf("-debugdir=%s/garble", g.ProjectDir))
	}

	goCommand = append(goCommand, []string{"-literals", "-tiny", "-seed=random"}...)

	goCommand = append(goCommand, command...)

	return goCommand
}

// IsTypeSupported retrieves build info based on type of executable
func (g *GoConfig) IsTypeSupported(t string) (string, string, error) {

	switch strings.ToLower(t) {
	case "exe":
		return "exe", "", nil
	case "dll":
		return "dll", "c-shared", nil
	case "svc":
		return "exe", "", nil
	case "cpl":
		return "cpl", "c-shared", nil
	case "xll":
		return "xll", "c-shared", nil
	case "pie":
		return "exe", "pie", nil
	default:
		return "", "", fmt.Errorf("Type must be exe, dll, cpl, xll, or pie.")
	}
}

// modified from https://github.com/BishopFox/sliver/blob/5bcfa4c249341e9c9032abcaaf1d4cf459e20059/server/gogo/go.go#L107
// GetGoCache - Get the OS temp dir (used for GOCACHE)
func getGoCache(appDir string) string {
	cachePath := path.Join(appDir, "cache")
	os.MkdirAll(cachePath, 0700)

	absPath, _ := filepath.Abs(cachePath)
	return absPath
}
