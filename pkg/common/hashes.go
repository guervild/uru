package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

func GetMD5Hash(payload []byte) string {
	return fmt.Sprintf("%x", md5.Sum(payload))
}

func GetSHA1Hash(payload []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(payload))
}

func GetSHA256Hash(payload []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(payload))
}
