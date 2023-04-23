package c

import (
	"embed"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

// taken from go implmentation
type CreateThread struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

// this is used specifically for cmd list
func NewCreateThread() models.ObjectModel {
	return &CreateThread{
		Name:        "windows/native/local/createthread",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use windows apis (virtual alloc, virtual protect, memcpy, CreateThread) to inject code",
	}
}

func (e *CreateThread) GetImports() []string {
	return []string{
		"windows.h",
	}
}

// taken from go implementation
func (e *CreateThread) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/c/injector/windows/native/local/createthread/instanciations.c.tmpl", e)
}

// taken from go implementation
func (e *CreateThread) RenderFunctionCode(data embed.FS) (string, error) {
	//return common.CommonRendering(data, "templates/c/injector/basicInjector_createThread/functions.c.tmpl", e)

	return "", nil
}
