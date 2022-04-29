package encoder

import (
	"embed"
	"encoding/hex"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type HexEncoder struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

func NewHexEncoder() models.ObjectModel {
	return &HexEncoder{
		Name:        "hex",
		Description: "Use hex encoding to encode given data",
		Debug:       false,
	}
}

func (e *HexEncoder) Encode(shellcode []byte) ([]byte, error) {

	return []byte(hex.EncodeToString(shellcode)), nil
}

func (e *HexEncoder) GetImports() []string {

	return []string{
		`"encoding/hex"`,
	}
}

func (e *HexEncoder) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/encoders/hex/instanciation.go.tmpl", e)
}

func (e *HexEncoder) RenderFunctionCode(data embed.FS) (string, error) {

	return "", nil
}
