package _go

import (
	"crypto/rc4"
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type RC4Encoder struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

func NewRC4Encoder() models.ObjectModel {
	return &RC4Encoder{
		Name:        "rc4",
		Description: "Use rc4 encoding to encode given data",
		Key:         common.RandomString(12),
		Debug:       false,
	}
}

func (e *RC4Encoder) Encode(shellcode []byte) ([]byte, error) {

	cipher, err := rc4.NewCipher([]byte(e.Key))

	if err != nil {
		return nil, err
	}

	encryptedBytes := make([]byte, len(shellcode))
	cipher.XORKeyStream(encryptedBytes, shellcode)

	return encryptedBytes, nil
}

func (e *RC4Encoder) GetImports() []string {

	return []string{
		`"crypto/rc4"`,
	}
}

func (e *RC4Encoder) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/rc4/instanciation.go.tmpl", e)
}

func (e *RC4Encoder) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/rc4/functions.go.tmpl", e)
}
