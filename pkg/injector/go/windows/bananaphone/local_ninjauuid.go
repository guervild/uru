package bananaphone

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type NinjaUUID struct {
	Name        string
	Description string
	Debug       bool
	Library     string
}

func NewNinjaUUID() models.ObjectModel {
	return &NinjaUUID{
		Name:        "windows/bananaphone/local/ninja_uuid",
		Description: "[experimental/dev] Module stomping following EnumSystemLocalesA for injection. Injection taken from @boku7 project. uuid encoder must be used as your last encoder.",
		Debug:       false,
		Library:     "C:\\\\windows\\\\system32\\\\windows.storage.dll",
	}
}

func (i *NinjaUUID) GetImports() []string {

	return []string{
		`"syscall"`,
		`"unsafe"`,
		`bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"`,
		`"golang.org/x/sys/windows"`,
	}
}

func (e *NinjaUUID) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/injector/windows/bananaphone/local/ninja_uuid/instanciation.go.tmpl", e)
}

func (e *NinjaUUID) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/injector/windows/bananaphone/local/ninja_uuid/functions.go.tmpl", e)
}
