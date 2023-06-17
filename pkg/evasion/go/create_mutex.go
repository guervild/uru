package _go

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type CreateMutexEvasion struct {
	Name        string
	Description string
	MutexName   string
	Debug       bool
}

func NewCreateMutexEvasion() models.ObjectModel {
	return &CreateMutexEvasion{
		Name: "CreateMutex",
		Description: `Create a mutex to have a single instance of the program running.
  Argument(s):
    MutexName: The mutex name. Default is "UruMutex".`,
		MutexName: "UruMutex",
		Debug:     false,
	}
}

func (e *CreateMutexEvasion) GetImports() []string {
	return []string{
		`"unsafe"`,
		`"syscall"`,
	}
}

func (e *CreateMutexEvasion) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/evasions/createmutex/instanciation.go.tmpl", e)
}

func (e *CreateMutexEvasion) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/evasions/createmutex/functions.go.tmpl", e)
}
