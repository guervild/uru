package encoder

import (
	"fmt"
	"github.com/guervild/uru/pkg/encoder/c"
	"github.com/guervild/uru/pkg/encoder/go"
	"github.com/guervild/uru/pkg/models"
)

func GetEncoder(encoderType string, langType string) (models.ObjectModel, error) {

	switch langType {
	case "go":
		if encoderType == "zip" {
			return _go.NewZipEncoder(), nil
		}

		if encoderType == "xor" {
			return _go.NewXorEncoder(), nil
		}

		if encoderType == "rc4" {
			return _go.NewRC4Encoder(), nil
		}

		if encoderType == "hex" {
			return _go.NewHexEncoder(), nil
		}

		if encoderType == "aes" {
			return _go.NewAESEncoder(), nil
		}

		if encoderType == "reverse_order" {
			return _go.NewReverseOrderEncoder(), nil
		}

		if encoderType == "uuid" {
			return _go.NewUUIDEncoder(), nil
		}

		if encoderType == "english_words" {
			return _go.NewEnglishWordsEncoder(), nil
		}

		return nil, fmt.Errorf("Wrong encoder type passed: encoder %s is unknown", encoderType)
	case "c":
		if encoderType == "xor" {
			return c.NewXorEncoder(), nil
		}
		return nil, fmt.Errorf("Wrong encoder type passed: encoder %s is unknown", encoderType)
	case "rust":
		return nil, fmt.Errorf("Wrong encoder type passed: encoder %s is unknown", encoderType)
	}
	return nil, fmt.Errorf("Wrong langtype: %s is unknown", langType)
}
