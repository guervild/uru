package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/encoder"
	"github.com/guervild/uru/pkg/evasion"
	"github.com/guervild/uru/pkg/injector"
	"github.com/guervild/uru/pkg/logger"

	"github.com/spf13/cobra"
)

var pkg_path = "./pkg/"

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

	// arg[0] = lang
	//  arg[1] = evasion, encoder etc.
	switch strings.ToLower(args[1]) {
	case "encoders":
		printValue("Encoders", listEncoders(strings.ToLower(args[0])))
	case "evasions":
		printValue("Evasions", listEvasions(strings.ToLower(args[0])))
	case "injectors":
		printValue("Injectors", listInjectors(strings.ToLower(args[0])))
	case "options":
		printValue("Options", getOptions())

	default:
		fmt.Println("Argument must be encoders, evasions, injectors or options.")
	}
}

func createListOfModules(module string, langType string) []string {

	// path to encoder go files (specified by lang)
	var dirPath = pkg_path + module + langType

	var lst = make([]string, 0)

	err := filepath.Walk(dirPath,
		func(path string, info fs.FileInfo, err error) error {

			if err != nil {
				return err
			}
			if !info.IsDir() {
				//string will be like:pkg/injector/c/basicInjector_executeFP.go

				// get rid of exstension
				currStr := strings.Split(path, ".go")[0]

				// get rid of pkg/[module]/[lang]
				currStr = strings.Split(currStr, "pkg/"+module+langType+"/")[1]
				lst = append(lst, currStr)
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

	encoders := createListOfModules("encoder/", langType)

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

	evasions := createListOfModules("evasion/", langType)

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

	injectors := createListOfModules("injector/", langType)

	for _, v := range injectors {

		injectorValue, err := injector.GetInjector(strings.ToLower(v), langType)
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
