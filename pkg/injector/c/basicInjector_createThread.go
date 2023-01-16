package c

import (
	"embed"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

// taken from go implmentation
type CBasicInjector_createThread struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

// this is used specifically for cmd list
func NewCBasicInjector_createThread() models.ObjectModel {
	return &CBasicInjector_createThread{
		Name:        "basicinjector_createthread",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use windows apis (virtual alloc, virtual protect, writeprocessmemory, executeFP) to inject code",
	}
}

func (e *CBasicInjector_createThread) GetImports() []string {
	return []string{
		"windows.h",
	}
}

// taken from go implementation
func (e *CBasicInjector_createThread) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/c/injector/basicInjector_createThread/instanciations.c.tmpl", e)
}

// taken from go implementation
func (e *CBasicInjector_createThread) RenderFunctionCode(data embed.FS) (string, error) {
	//return common.CommonRendering(data, "templates/c/injector/basicInjector_createThread/functions.c.tmpl", e)

	return "", nil
}
