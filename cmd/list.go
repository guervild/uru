package cmd

import (
	"fmt"
	"io/fs"
	"strings"
	"path/filepath"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/encoder"
	"github.com/guervild/uru/pkg/evasion"
	"github.com/guervild/uru/pkg/injector"
	"github.com/guervild/uru/pkg/logger"
	"github.com/guervild/uru/data"

	"github.com/spf13/cobra"
)


var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List basic options, and artifacts (encoders, evasions and injectors).",
	Long:  `List basic options, and artifacts (encoders, evasions and injectors). Accept : "encoders", "evasions", "injectors", or "options"`,
	Args:  cobra.RangeArgs(1, 2),
	Run:   List,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List(cmd *cobra.Command, args []string) {


	if strings.ToLower(args[0]) != "options" {
		if len(args) != 2 { 
			fmt.Println("Need a second argument 'go' or 'c'")
			return
		}

		if strings.ToLower(args[1]) != "go" &&  strings.ToLower(args[1]) != "c" {
			fmt.Println("Language must be'go' or 'c'")
			return
		}
	}

	switch strings.ToLower(args[0]) {
	case "encoders":
		printValue("Encoders", listEncoders(strings.ToLower(args[1])))
	case "evasions":
		printValue("Evasions", listEvasions(strings.ToLower(args[1])))
	case "injectors":
		printValue("Injectors", listInjectors(strings.ToLower(args[1])))
	case "options":
		printValue("Options", getOptions())

	default:
		fmt.Println("Argument must be encoders, evasions, injectors or options. For encoders, evasions and injectors, a second argument can be passed 'go' or 'c'.")
	}
}

func createListOfModules(module string, langType string) []string {

	dataTmpl := data.GetTemplates()
	pathLang := filepath.Join("./templates", langType)
	pathModule := filepath.Join(pathLang, module)

	var lst = make([]string, 0)

	err := fs.WalkDir(dataTmpl, pathModule,
		 func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				logger.Logger.Fatal().Msg(err.Error())
			}

			if d.IsDir() {
				if strings.ToLower(d.Name()) != "commons" && strings.ToLower(d.Name()) != strings.ToLower(module) {
					if strings.ToLower(module) == "injector" {
						lst = append(lst, path)
					} else {
						lst = append(lst, d.Name())
					}
				}
			}

			return nil
		})

	if err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}
	return lst
}

func listEncoders(langType string) map[string]string {

	mEncoders := make(map[string]string)

	encoders := createListOfModules("encoders", langType)

	for _, v := range encoders {

		encoderValue, err := encoder.GetEncoder(strings.ToLower(v), langType)
		if err != nil {
			logger.Logger.Info().Msg(err.Error())
		} else {
			name := common.GetField(encoderValue, "Name")
			desc := common.GetField(encoderValue, "Description")

			mEncoders[name] = desc
		}
	}

	return mEncoders
}

func listEvasions(langType string) map[string]string {
	mEvasions := make(map[string]string)

	evasions := createListOfModules("evasions", langType)

	for _, v := range evasions {

		evasionValue, err := evasion.GetEvasion(strings.ToLower(v), langType)
		if err != nil {
			// logger.Logger.Info().Msg(err.Error())
		} else {
			name := common.GetField(evasionValue, "Name")
			desc := common.GetField(evasionValue, "Description")

			mEvasions[name] = desc
		}
	}
	return mEvasions
}

func listInjectors(langType string) map[string]string {
	mInjectors := make(map[string]string)

	injectors := createListOfModules("injector", langType)

	for _, v := range injectors {
		path := strings.ReplaceAll(v, fmt.Sprintf("templates/%s/injector/", langType), "")
		//path = strings.ReplaceAll(path, "local/", "")
		//path = strings.ReplaceAll(path, "-", "")

		injectorValue, err := injector.GetInjector(strings.ToLower(path), langType)
		if err != nil {
			logger.Logger.Info().Msg(err.Error())
		} else {
			name := common.GetField(injectorValue, "Name")
			desc := common.GetField(injectorValue, "Description")

			mInjectors[name] = desc
		}
	}
	return mInjectors
}

func printValue(title string, value map[string]string) {

	fmt.Println(fmt.Sprintf("========== %s ==========\n", title))

	for k, v := range value {
		n := fmt.Sprintf("  Name       : %s", k)
		d := fmt.Sprintf("  Description: %s\n", v)
		fmt.Println(n)
		fmt.Println(d)
	}
}

func getOptions() map[string]string {

	mOptions := make(map[string]string)

	mOptions["arch"] = "Architecture of the compiled program: x64 or x86. (mandatory)"
	mOptions["type"] = "Type of payload to compiled. Can be exe, pie or dll. (mandatory)"
	//[SGN] - DECOMMENT TO USE SGN
	//mOptions["sgn"] = "Apply SGN on the provided payload file (might not work with all shellcode/executable)."
	mOptions["debug"] = "Add print and debug functions to the compiled program."
	mOptions["obfuscation"] = "Obfuscate the code before compilation (Use garble to obfuscate the code)."
	mOptions["append"] = "Append the provided string at the end of the payload file. Must be hexadecimal ex: 90909090"
	mOptions["preprend"] = "Prepend the payload file with the provided string. Must be hexadecimal ex: 90909090"

	return mOptions
}
