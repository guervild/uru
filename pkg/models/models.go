package models

import (
	"embed"
)

type ObjectModel interface {
	GetImports() []string
	RenderInstanciationCode(data embed.FS) (string, error)
	RenderFunctionCode(data embed.FS) (string, error)
}
