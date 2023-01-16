package _go

import (
	"embed"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type RefreshDllEvasion struct {
	Name         string
	Description  string
	DllName      string
	UseBanana    string
	SubNameError string
	Debug        bool
}

func NewRefreshDllEvasion() models.ObjectModel {
	return &RefreshDllEvasion{
		Name: "RefreshDll",
		Description: `Refresh the given dll to remove hook by using the dll on disk. (Inspired by sliver/scarecrow and TomWhitez works).
  Argument(s):
    UseBanana: UseBananaPhone to perform syscall. Default is "false".
	DllName: Name of the dll to refresh. Default is "C:\\\\Windows\\\\System32\\\\kernel32.dll".`,
		UseBanana:    "false",
		DllName:      "C:\\\\Windows\\\\System32\\\\kernel32.dll",
		Debug:        false,
		SubNameError: common.RandomStringOnlyChar(5),
	}
}

func (e *RefreshDllEvasion) GetImports() []string {

	imports := []string{
		`"debug/pe"`,
		`"io/ioutil"`,
		`"fmt"`,
		`"unsafe"`,
		`"golang.org/x/sys/windows"`,
	}

	if strings.ToLower(e.UseBanana) == "true" {
		imports = append(imports, `bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"`)
	}

	return imports
}

func (e *RefreshDllEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/refreshdll/instanciation.go.tmpl", e)
}

func (e *RefreshDllEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/refreshdll/functions.go.tmpl", e)
}
