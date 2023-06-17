package native

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type NtQueueApcThreadExLocal struct {
	Name        string
	Description string
	Debug       bool
}

func NewNtQueueApcThreadExLocal() models.ObjectModel {
	return &NtQueueApcThreadExLocal{
		Name:        "windows/native/local/ntqueueapcthreadex",
		Description: "Use native windows api call NtQueueApcThreadEx to inject in the current process",
		Debug:       false,
	}
}

func (i *NtQueueApcThreadExLocal) GetImports() []string {
	return []string{
		`"syscall"`,
		`"unsafe"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (i *NtQueueApcThreadExLocal) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/injector/windows/native/local/ntqueueapcthreadex/instanciation.go.tmpl", i)
}

func (i *NtQueueApcThreadExLocal) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/injector/windows/native/local/ntqueueapcthreadex/functions.go.tmpl", i)
}
