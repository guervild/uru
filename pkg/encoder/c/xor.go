package c

import (
	"embed"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

// taken from go implmentation
type xorEncoder struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

// this is used specifically for cmd list
func NewXorEncoder() models.ObjectModel {
	return &xorEncoder{
		Name:        "xor",
		Key:         common.RandomStringOnlyChar(12),
		Description: "Use xor algorithm to encode given data",
	}
}

func (e *xorEncoder) Encode(shellcode []byte) ([]byte, error) {
	var output []byte

	kL := len(e.Key)

	for i := range shellcode {
		output = append(output, byte(shellcode[i]^e.Key[i%kL]))
	}

	return output, nil
}

func (e *xorEncoder) GetImports() []string {
	return []string{
		"windows.h",
	}
}

// taken from go implementation
func (e *xorEncoder) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/c/encoders/xor/instanciations.c.tmpl", e)
}

// taken from go implementation
func (e *xorEncoder) RenderFunctionCode(data embed.FS) (string, error) {
	//return common.CommonRendering(data, "templates/c/encoders/xor/functions.c.tmpl", e)
	return "", nil
}
