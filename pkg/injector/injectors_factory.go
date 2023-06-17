package injector

import (
	"fmt"

	cnative "github.com/guervild/uru/pkg/injector/c/windows/native"
	gobananaphone "github.com/guervild/uru/pkg/injector/go/windows/bananaphone"
	gonative "github.com/guervild/uru/pkg/injector/go/windows/native"

	"github.com/guervild/uru/pkg/models"
)

func GetInjector(injectorType string, langType string) (models.ObjectModel, error) {
	switch langType {
	case "go":
		if injectorType == "windows/native/local/execute_fp" {
			return gonative.NewExecuteFP(), nil
		}
		if injectorType == "windows/native/local/ntqueueapcthreadex" {
			return gonative.NewNtQueueApcThreadExLocal(), nil
		}
		if injectorType == "windows/native/local/createthread" {
			return gonative.NewCreateThreadNative(), nil
		}
		if injectorType == "windows/bananaphone/local/ntqueueapcthreadex" {
			return gobananaphone.NewNtQueueApcThreadExLocal(), nil
		}
		if injectorType == "windows/bananaphone/local/execute_fp" {
			return gobananaphone.NewExecuteFP(), nil
		}
		if injectorType == "windows/bananaphone/local/ninja_uuid" {
			return gobananaphone.NewNinjaUUID(), nil
		}
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	case "c":
		switch injectorType {
		case "windows/native/local/createthread":
			return cnative.NewCreateThread(), nil
		case "windows/native/local/execute_fp":
			return cnative.NewExecuteFP(), nil
		default:
			break
		}
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	case "rust":
		return nil, fmt.Errorf("Wrong injector type passed: injector %s is unknown", injectorType)
	}
	return nil, fmt.Errorf("Wrong langtype: %s is unknown", langType)
}
