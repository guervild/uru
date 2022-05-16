package cmd

import (
	"fmt"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/encoder"
	"github.com/guervild/uru/pkg/evasion"
	"github.com/guervild/uru/pkg/injector"
	"github.com/guervild/uru/pkg/logger"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List basic options, and artifacts (encoders, evasions and injectors).",
	Long:  `List basic options, and artifacts (encoders, evasions and injectors). Accept : "encoders", "evasions", "injectors", or "options"`,
	Args:  cobra.ExactArgs(1),
	Run:   List,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List(cmd *cobra.Command, args []string) {

	switch strings.ToLower(args[0]) {
	case "encoders":
		printValue("Encoders", listEncoders())
	case "evasions":
		printValue("Evasions", listEvasions())
	case "injectors":
		printValue("Injectors", listInjectors())
	case "options":
		printValue("Options", getOptions())

	default:
		fmt.Println("Argument must be encoders, evasions, injectors or options.")
	}
}

func listEncoders() map[string]string {

	mEncoders := make(map[string]string)

	encoders := []string{
		"xor",
		"zip",
		"rc4",
		"hex",
		"aes",
		"reverse-order",
		"uuid",
	}

	for _, v := range encoders {

		encoderValue, err := encoder.GetEncoder(strings.ToLower(v))
		if err != nil {
			logger.Logger.Fatal().Msg(err.Error())
		}

		name := common.GetField(encoderValue, "Name")
		desc := common.GetField(encoderValue, "Description")

		mEncoders[name] = desc
	}

	return mEncoders
}

func listEvasions() map[string]string {
	mEvasions := make(map[string]string)

	evasions := []string{
		"sleep",
		"hideconsole",
		"isdomainjoined",
		"selfdelete",
		"ntsleep",
		"english-words",
		"patchamsi",
		"patchetw",
		"patch",
		"createmutex",
	}

	for _, v := range evasions {

		evasionValue, err := evasion.GetEvasion(strings.ToLower(v))
		if err != nil {
			logger.Logger.Fatal().Msg(err.Error())
		}

		name := common.GetField(evasionValue, "Name")
		desc := common.GetField(evasionValue, "Description")

		mEvasions[name] = desc
	}

	return mEvasions
}

func listInjectors() map[string]string {
	mInjectors := make(map[string]string)

	injectors := []string{
		"windows/native/local/CreateThreadNative",
		"windows/native/local/ntqueueapcthreadex-local",
		"windows/native/local/go-shellcode-syscall",
		"windows/bananaphone/local/ntqueueapcthreadex-local",
		"windows/bananaphone/local/go-shellcode-syscall",
		"windows/bananaphone/local/ninjauuid",
	}

	for _, v := range injectors {

		injectorValue, err := injector.GetInjector(strings.ToLower(v))
		if err != nil {
			logger.Logger.Fatal().Msg(err.Error())
		}

		name := common.GetField(injectorValue, "Name")
		desc := common.GetField(injectorValue, "Description")

		mInjectors[name] = desc
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
