package injector

import (
	"fmt"
	"github.com/guervild/uru/pkg/injector/c"
	bananaphone2 "github.com/guervild/uru/pkg/injector/go/windows/bananaphone"
	native2 "github.com/guervild/uru/pkg/injector/go/windows/native"

	"github.com/guervild/uru/pkg/models"
)

func GetInjector(injectorType string, langType string) (models.ObjectModel, error) {

	switch langType {
	case "go":
		if injectorType == "windows/native/syscall" {
			return native2.NewSyscallGoShellcode(), nil
		}
		if injectorType == "windows/native/ntqueueapcthreadexlocal" {
			return native2.NewNtQueueApcThreadExLocal(), nil
		}
		if injectorType == "windows/native/createthreadnative" {
			return native2.NewCreateThreadNative(), nil
		}
		if injectorType == "windows/bananaphone/ntqueueapcthreadexlocal" {
			return bananaphone2.NewNtQueueApcThreadExLocal(), nil
		}
		if injectorType == "windows/bananaphone/syscall" {
			return bananaphone2.NewSyscallGoShellcode(), nil
		}
		if injectorType == "windows/bananaphone/ninjauuid" {
			return bananaphone2.NewNinjaUUID(), nil
		}
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	case "c":
		switch injectorType {
		case "basicinjector_createthread":
			return c.NewCBasicInjector_createThread(), nil
		case "basicinjector_executefp":
			return c.NewCBasicInjector_executeFP(), nil
		default:
			break
		}
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	case "rust":
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	}
	return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
}
