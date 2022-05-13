package encoder

import (
	"io"
	"bytes"
	"strings"

	"github.com/Binject/debug/pe"
	"github.com/Binject/go-donut/donut"
)

func DetectDotNet(r io.ReaderAt) (bool, string) {
	// auto-detect .NET assemblies and version
	pefile, err := pe.NewFile(r)
	if err != nil {
		return false, ""
	}
	defer pefile.Close()
	return pefile.IsManaged(), pefile.NetCLRVersion()
}

func ConvertToGoDonutShellcode(payload []byte, extension, class, method, parameters string) ([]byte, error){
	reader := bytes.NewReader(payload)
	var dType donut.ModuleType
	var runtime string

	switch strings.ToLower(extension) {
	case ".exe":
		dotNetMode, dotNetVersion := DetectDotNet(reader)
		if dotNetMode {
			dType = donut.DONUT_MODULE_NET_EXE
		} else {
			dType = donut.DONUT_MODULE_EXE
		}
		if dotNetVersion != "" && runtime == "" {
			runtime = dotNetVersion
		}
	case ".dll":
		dotNetMode, dotNetVersion := DetectDotNet(reader)
		if dotNetMode {
			dType = donut.DONUT_MODULE_NET_DLL
		} else {
			dType = donut.DONUT_MODULE_DLL
		}
		if dotNetVersion != "" && runtime == "" {
			runtime = dotNetVersion
		}
	case ".xsl":
		dType = donut.DONUT_MODULE_XSL
	case ".js":
		dType = donut.DONUT_MODULE_JS
	case ".vbs":
		dType = donut.DONUT_MODULE_VBS
	}
	
	shellcode, err := donut.ShellcodeFromBytes(bytes.NewBuffer(payload), &donut.DonutConfig{
		Arch:       donut.X84,
		Type:       dType,
		InstType:   donut.DONUT_INSTANCE_PIC,
		Entropy:    donut.DONUT_ENTROPY_DEFAULT,
		Runtime:	runtime,
		Class: class,
		Method: method,
		Compress:   1,
		Format:     1,
		Bypass:     3,
		Parameters: parameters,
	})

	return shellcode.Bytes(), err
}