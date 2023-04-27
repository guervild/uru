package c

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type CSleepEvasion struct {
	Name        string
	Debug       bool
	Delay       string
	Description string
}

func NewCSleepEvasion() models.ObjectModel {
	return &CSleepEvasion{
		Name:  "sleep",
		Debug: false,
		Delay: "5000",
		Description: `Sleep during a fixed amount of time in seconds.
  Argument(s):
    Delay: default value is 5s`,
	}
}

func (e *CSleepEvasion) GetImports() []string {

	return []string{
		"windows.h",
	}
}

func (e *CSleepEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/c/evasions/sleep/instanciations.c.tmpl", e)
}

func (e *CSleepEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return "", nil
}
