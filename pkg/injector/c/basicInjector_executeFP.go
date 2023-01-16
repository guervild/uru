package c

import (
	"embed"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

// taken from go implmentation
type CBasicInjector_executeFP struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

// this is used specifically for cmd list
func NewCBasicInjector_executeFP() models.ObjectModel {
	return &CBasicInjector_executeFP{
		Name:        "basicinjector_executefp",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use windows apis (virtual alloc, virtual protect, writeprocessmemory, executeFP) to inject code",
	}
}

func (e *CBasicInjector_executeFP) GetImports() []string {
	return []string{
		"windows.h",
	}
}

// taken from go implementation
func (e *CBasicInjector_executeFP) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/c/injector/basicInjector_executeFP/instanciations.c.tmpl", e)
}

// taken from go implementation
func (e *CBasicInjector_executeFP) RenderFunctionCode(data embed.FS) (string, error) {
	return "", nil
}
