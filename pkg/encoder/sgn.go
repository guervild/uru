package encoder

//[SGN] - DECOMMENT TO USE SGN
/*
import (
	"fmt"

	sgn "github.com/EgeBalci/sgn/pkg"
)

func DoSGNEncode(architecture string, payload []byte) ([]byte, error) {

	arch, err := getProperArch(architecture)
	if err != nil {
		return nil, err
	}
	// Create a new SGN encoder
	encoder := sgn.NewEncoder()
	// Set the proper architecture
	encoder.SetArchitecture(arch)
	// Encode the binary
	encodedBinary, err := encoder.Encode(payload)

	return encodedBinary, err
}

func getProperArch(architecture string) (int, error) {
	if architecture == "x86" {
		return 86,nil
	}

	if architecture == "x64" {
		return 64,nil
	}

	return 0, fmt.Errorf("Architecture for SGN is invalid, must be x64 or x86: %s", architecture)
}*/
