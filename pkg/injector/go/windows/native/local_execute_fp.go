package native

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type ExecuteFP struct {
	Name        string
	Description string
	Debug       bool
}

func NewExecuteFP() models.ObjectModel {
	return &ExecuteFP{
		Name:        "windows/native/local/execute_fp",
		Description: "Executes Shellcode in the current running proccess by making a Syscall on the Shellcode's entry point.",
		Debug:       false,
	}
}

func (i *ExecuteFP) GetImports() []string {
	return []string{
		`"syscall"`,
		`"unsafe"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (i *ExecuteFP) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/injector/windows/native/local/execute_fp/instanciation.go.tmpl", i)
}

func (i *ExecuteFP) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/injector/windows/native/local/execute_fp/functions.go.tmpl", i)
}
