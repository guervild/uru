package bananaphone

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
		Name:        "windows/bananaphone/local/execute_fp",
		Description: "Executes Shellcode in the current running proccess by making a Syscall on the Shellcode's entry point. Syscall (memory allocation) is performed using using bananaphone from @C-Sto.",
		Debug:       false,
	}
}

func (i *SyscallGoShellcode) GetImports() []string {

	return []string{
		`"syscall"`,
		`"unsafe"`,
		`bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"`,
	}
}

func (e *SyscallGoShellcode) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/injector/windows/bananaphone/local/go-shellcode-syscall/instanciation.go.tmpl", e)
}

func (e *SyscallGoShellcode) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/injector/windows/bananaphone/local/go-shellcode-syscall/functions.go.tmpl", e)
}
