package native

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type SyscallGoShellcode struct {
	Name        string
	Description string
	Debug       bool
}

func NewSyscallGoShellcode() models.ObjectModel {
	return &SyscallGoShellcode{
		Name:        "windows/native/local/go-shellcode-syscall",
		Description: "Executes Shellcode in the current running proccess by making a Syscall on the Shellcode's entry point.",
		Debug:       false,
	}
}

func (i *SyscallGoShellcode) GetImports() []string {

	return []string{
		`"syscall"`,
		`"unsafe"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (e *SyscallGoShellcode) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/injector/windows/native/local/go-shellcode-syscall/instanciation.go.tmpl", e)
}

func (e *SyscallGoShellcode) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/injector/windows/native/local/go-shellcode-syscall/functions.go.tmpl", e)
}
