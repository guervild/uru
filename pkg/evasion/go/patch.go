package _go

import (
	"embed"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type PatchEvasion struct {
	Name        string
	Description string
	Debug       bool
	Module      string
	Proc        string
	Data        string
	UseBanana   string
	SubName     string
}

func NewPatchEvasion() models.ObjectModel {

	return &PatchEvasion{
		Name: "patch",
		Description: `Path a given method (By default patch back EtwEventWrite). Can use BananaPhone if set. (credits: method taken from Merlin, @Ne0nd0g)
  Argument(s):
    Module: the module where the function is. Example: "ntdll.dll"
    Proc: the function to patch. Example: "EtwEventWrite"
    Data: the data to use to patch the function in hex.
    UseBanana: use BananaPhone to perfom syscall`,
		Module:    "ntdll.dll",
		Proc:      "EtwEventWrite",
		Data:      "4C8BDC48",
		UseBanana: "false",
		SubName:   common.RandomStringOnlyChar(5),
	}
}

func (e *PatchEvasion) GetImports() []string {

	imports := []string{
		`"encoding/hex"`,
		`"fmt"`,
		`"syscall"`,
		`"unsafe"`,
	}

	if strings.ToLower(e.UseBanana) == "true" {
		imports = append(imports, `bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"`)
	} else {
		imports = append(imports, `"golang.org/x/sys/windows"`)
	}

	return imports
}

func (e *PatchEvasion) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/evasions/patch/instanciation.go.tmpl", e)
}

func (e *PatchEvasion) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/evasions/patch/patch_functions.go.tmpl", e)
}
