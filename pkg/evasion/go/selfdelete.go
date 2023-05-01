package _go

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type SelfDeleteEvasion struct {
	Name        string
	Description string
	NewADS      string
	Debug       bool
}

func NewSelfDeleteEvasion() models.ObjectModel {
	return &SelfDeleteEvasion{
		Name:        "selfdelete",
		Description: "Delete the current binary after shellcode/payload execution",
		NewADS:      common.RandomStringOnlyChar(4),
		Debug:       false,
	}
}

func (e *SelfDeleteEvasion) GetImports() []string {

	return []string{
		`"syscall"`,
		`"unsafe"`,
		`"fmt"`,
		`"unicode/utf16"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (e *SelfDeleteEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/selfdelete/instanciation.go.tmpl", e)
}

func (e *SelfDeleteEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/selfdelete/functions.go.tmpl", e)
}
