package _go

//inspired by https://github.com/D00MFist/Go4aRun/blob/493acbb0c38be9719dbfa7f44e5d4d9d144709c9/pkg/useful/useful.go#L27
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"embed"
	"io"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type AESEncoder struct {
	Name        string
	Description string
	Key         string
	Debug       bool
}

func NewAESEncoder() models.ObjectModel {
	return &AESEncoder{
		Name:        "AES",
		Description: "Use AES GCM to encrypt given data",
		Key:         common.RandomString(32),
		Debug:       false,
	}
}

func (e *AESEncoder) Encode(shellcode []byte) ([]byte, error) {

	block, _ := aes.NewCipher([]byte(e.Key))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, shellcode, nil)

	return ciphertext, nil
}

func (e *AESEncoder) GetImports() []string {

	return []string{
		`"crypto/aes"`,
		`"crypto/cipher"`,
	}
}

func (e *AESEncoder) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/aes/instanciation.go.tmpl", e)
}

func (e *AESEncoder) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/aes/functions.go.tmpl", e)
}
