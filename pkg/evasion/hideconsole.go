package evasion

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type HideConsoleEvasion struct {
	Name        string
	Description string
}

func NewHideConsoleEvasion() models.ObjectModel {
	return &HideConsoleEvasion{
		Name:        "hideconsole",
		Description: "Prevent windows console to be displayed",
	}
}

func (e *HideConsoleEvasion) GetImports() []string {

	return []string{
		`"syscall"`,
	}
}

func (e *HideConsoleEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/evasions/hideconsole/instanciation.go.tmpl", e)
}

func (e *HideConsoleEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/evasions/hideconsole/functions.go.tmpl", e)
}
