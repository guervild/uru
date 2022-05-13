package evasion

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type NtSleepEvasion struct {
	Name        string
	Delay       string
	Debug       bool
	Description string
}

func NewNtSleepEvasion() models.ObjectModel {
	return &NtSleepEvasion{
		Name:        "NtSleep",
		Debug:       false,
		Delay:       "5",
		Description: `NtSleep during a fixed amount of time in seconds using NtDelayExecution API call
  Argument(s):
    Delay: default value is 5s`,
	}
}

func (e *NtSleepEvasion) GetImports() []string {

	return []string{
		`"unsafe"`,
		`"syscall"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (e *NtSleepEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/evasions/ntsleep/instanciation.go.tmpl", e)
}

func (e *NtSleepEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/evasions/ntsleep/functions.go.tmpl", e)
}
