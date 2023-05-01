package _go

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type SleepEvasion struct {
	Name        string
	Debug       bool
	Delay       string
	Description string
}

func NewSleepEvasion() models.ObjectModel {
	return &SleepEvasion{
		Name:  "sleep",
		Debug: false,
		Delay: "5",
		Description: `Sleep during a fixed amount of time in seconds.
  Argument(s):
    Delay: default value is 5s`,
	}
}

func (e *SleepEvasion) GetImports() []string {

	return []string{
		`"time"`,
	}
}

func (e *SleepEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/sleep/instanciation.go.tmpl", e)
}

func (e *SleepEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return "", nil
}
