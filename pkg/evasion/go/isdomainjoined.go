package _go

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type IsDomainJoinedEvasion struct {
	Name        string
	Description string
}

func NewIsDomainJoinedEvasion() models.ObjectModel {
	return &IsDomainJoinedEvasion{
		Name:        "IsDomainJoined",
		Description: "check if current computer is joined to a domain.",
	}
}

func (e *IsDomainJoinedEvasion) GetImports() []string {

	return []string{
		`"syscall"`,
	}
}

func (e *IsDomainJoinedEvasion) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/isdomainjoined/instanciation.go.tmpl", e)
}

func (e *IsDomainJoinedEvasion) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/evasions/isdomainjoined/functions.go.tmpl", e)
}
