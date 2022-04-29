package encoder

import (
	"fmt"
	"github.com/guervild/uru/pkg/models"
)

func GetEncoder(encoderType string) (models.ObjectModel, error) {
	if encoderType == "zip" {
		return NewZipEncoder(), nil
	}

	if encoderType == "xor" {
		return NewXorEncoder(), nil
	}

	if encoderType == "rc4" {
		return NewRC4Encoder(), nil
	}

	if encoderType == "hex" {
		return NewHexEncoder(), nil
	}

	if encoderType == "aes" {
		return NewAESEncoder(), nil
	}

	if encoderType == "reverse-order" {
		return NewReverseOrderEncoder(), nil
	}

	if encoderType == "uuid" {
		return NewUUIDEncoder(), nil
	}

	return nil, fmt.Errorf("Wrong encoder type passed: encoder %s is unknown", encoderType)
}
