package _go

import (
	"archive/zip"
	"bytes"
	"embed"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type ZipEncoder struct {
	Name        string
	Description string
	Debug       bool
}

func NewZipEncoder() models.ObjectModel {
	return &ZipEncoder{
		Name:        "zip",
		Description: "Use zip compression on given data",
		Debug:       false,
	}
}

func (e *ZipEncoder) Encode(shellcode []byte) ([]byte, error) {

	var b bytes.Buffer

	z := zip.NewWriter(&b)

	f, err := z.Create(common.RandomString(12))

	if err != nil {
		return nil, err
	}

	if _, err := f.Write(shellcode); err != nil {
		return nil, err
	}

	if err := z.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *ZipEncoder) GetImports() []string {

	return []string{
		`"archive/zip"`,
		`"bytes"`,
		`"io/ioutil"`,
	}
}

func (e *ZipEncoder) RenderInstanciationCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/zip/instanciation.go.tmpl", e)
}

func (e *ZipEncoder) RenderFunctionCode(data embed.FS) (string, error) {

	return common.CommonRendering(data, "templates/go/encoders/zip/functions.go.tmpl", e)
}
