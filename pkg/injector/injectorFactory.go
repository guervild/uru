package injector

import (
	"fmt"

	"github.com/guervild/uru/pkg/injector/windows/bananaphone"
	"github.com/guervild/uru/pkg/injector/windows/native"
	"github.com/guervild/uru/pkg/models"
)

func GetInjector(injectorType string) (models.ObjectModel, error) {
	if injectorType == "windows/native/local/go-shellcode-syscall" {
		return native.NewSyscallGoShellcode(), nil
	}

	if injectorType == "windows/native/local/ntqueueapcthreadex-local" {
		return native.NewNtQueueApcThreadExLocal(), nil
	}

	if injectorType == "windows/native/local/createthreadnative" {
		return native.NewCreateThreadNative(), nil
	}

	if injectorType == "windows/bananaphone/local/ntqueueapcthreadex-local" {
		return bananaphone.NewNtQueueApcThreadExLocal(), nil
	}

	if injectorType == "windows/bananaphone/local/go-shellcode-syscall" {
		return bananaphone.NewSyscallGoShellcode(), nil
	}

	if injectorType == "windows/bananaphone/local/ninjauuid" {
		return bananaphone.NewNinjaUUID(), nil
	}

	return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
}
