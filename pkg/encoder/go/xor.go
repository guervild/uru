package _go

import (
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type XorEncoder struct {
	Name        string
	Key         string
	Description string
	Debug       bool
}

func NewXorEncoder() models.ObjectModel {
	return &XorEncoder{
		Name:        "xor",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use xor algorithm to encode given data",
	}
}

func (e *XorEncoder) Encode(shellcode []byte) ([]byte, error) {
	var output []byte

	kL := len(e.Key)

	for i := range shellcode {
		output = append(output, shellcode[i]^e.Key[i%kL])
	}

	return output, nil
}

func (e *XorEncoder) GetImports() []string {
	return nil
}

func (e *XorEncoder) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/encoders/xor/instanciation.go.tmpl", e)
}

func (e *XorEncoder) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/encoders/xor/functions.go.tmpl", e)
}
