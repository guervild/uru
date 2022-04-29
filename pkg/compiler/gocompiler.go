// modified from https://github.com/BishopFox/sliver/blob/5bcfa4c249341e9c9032abcaaf1d4cf459e20059/server/gogo/go.go

package compiler

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

func NewGoConfig(targetOs, arch, dirPath, buildmode string, keep, trimpath, obfuscation bool, imports []string) *GoConfig {

	config := &GoConfig{
		ProjectDir:  dirPath,
		GOOS:        targetOs,
		GOARCH:      arch,
		Buildmode:   buildmode,
		Trimpath:    trimpath,
		Obfuscation: obfuscation,
		Keep:        keep,
		Imports:	 imports,
	}

	config.GOCACHE = getGoCache(dirPath)

	return config
}

func (g *GoConfig) GoBuild(payload, dest string) ([]byte, error) {

	g.BuildGoBuildCommand(payload, dest)

	if 0 >= len(g.Command) {
		return nil, fmt.Errorf("Error while setting golang commands ...")
	}

	logger.Logger.Debug().Str("compile_args", strings.Join(g.Command," ")).Msg("Defining compile arguments")

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

	if g.Obfuscation {
		compilerPath = path.Join(path.Join(gopath, "bin"), "garble")
	} else {
		compilerPath = path.Join(path.Join(goroot, "bin"), "go")
	}

	g.CompilerPath = compilerPath

	logger.Logger.Debug().Str("compiler_path", g.CompilerPath).Msg("Set the compiler path")

	//Getting the right env variables
	env := []string{fmt.Sprintf("GOCACHE=%s", g.GOCACHE), fmt.Sprintf("GOPATH=%s", gopath), fmt.Sprintf("GOOS=%s", g.GOOS), fmt.Sprintf("GOARCH=%s", g.GOARCH), fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
	if g.Buildmode == "c-shared" {
		env = append(env, "CGO_ENABLED=1")
		env = append(env, "CXX=x86_64-w64-mingw32-g++")
		env = append(env, "CC=x86_64-w64-mingw32-gcc")
	}
	var gogarble []string
	if g.Obfuscation {
		for _,v := range g.Imports {
			var out string
			out = strings.TrimLeft(strings.TrimRight(v,"\""),"\"")
			if strings.Contains(out, "\"") {
				out2 := strings.Join(strings.Split(out, "\"")[1:], "\"")
				out = out2
			}
			gogarble = append(gogarble, out)
		}
		env = append(env, fmt.Sprintf("GOGARBLE=%s", strings.Join(gogarble, ",")))
	}

	g.Env = env

	logger.Logger.Debug().Str("compile_env_vars", strings.Join(g.Env," ")).Msg("Set the environment variables")

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

//modified from https://github.com/BishopFox/sliver/blob/5bcfa4c249341e9c9032abcaaf1d4cf459e20059/server/gogo/go.go#L107
// GetGoCache - Get the OS temp dir (used for GOCACHE)
func getGoCache(appDir string) string {
	cachePath := path.Join(appDir, "cache")
	os.MkdirAll(cachePath, 0700)

	absPath, _ := filepath.Abs(cachePath)
	return absPath
}