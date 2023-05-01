package c

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

// taken from go implmentation
type ExecuteFP struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

// this is used specifically for cmd list
func NewExecuteFP() models.ObjectModel {
	return &ExecuteFP{
		Name:        "windows/native/local/execute_fp",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use windows apis (virtual alloc, virtual protect, memcpy) to inject code",
	}
}

func (e *ExecuteFP) GetImports() []string {
	return []string{
		"windows.h",
	}
}

// taken from go implementation
func (e *ExecuteFP) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/c/injector/windows/native/local/execute_fp/instanciations.c.tmpl", e)
}

// taken from go implementation
func (e *ExecuteFP) RenderFunctionCode(data embed.FS) (string, error) {
	return "", nil
}
