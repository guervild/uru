package data

import (
	"embed"
)

//go:embed templates
var templates embed.FS

func GetTemplates() embed.FS {

	return templates
}
