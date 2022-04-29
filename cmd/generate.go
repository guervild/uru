package cmd

import (
	"io/ioutil"
	"path/filepath"
	"os"

	"github.com/guervild/uru/pkg/builder"
	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/logger"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a payload",
	Long:  `Take as input a config and a given payload file to generate a payload.`,
	Run:   Generate,
}
var (
	Payload    string
	Config     string
	Executable bool
	//Keep       bool
	Parameters string
	Output     string
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&Payload, "payload", "p", "", "`Shellcode/Executable` to use in the generated payload")
	generateCmd.MarkFlagRequired("payload")
	generateCmd.Flags().StringVarP(&Config, "config", "c", "", "Config file that definied the modules to use")
	generateCmd.MarkFlagRequired("config")
	generateCmd.Flags().StringVarP(&Parameters, "parameters", "", "", "Parameters to pass to the executable")
	generateCmd.Flags().StringVarP(&Output, "output", "o", "", "Output file name")
	generateCmd.Flags().BoolVarP(&Executable, "exe", "", false, "Process the given payload as an executable")
	//generateCmd.Flags().BoolVarP(&Executable, "keep", "", false, "Keep the content of the out directory (generated code, but also obfuscated code and cache if obfuscation is set to true)")

}

func Generate(cmd *cobra.Command, args []string) {

	//Check files
	if err := common.CheckIfFileExists(Payload); err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}

	if err := common.CheckIfFileExists(Config); err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}

	//Process config file
	configFilepath, _ := filepath.Abs(Config)
	configData, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}

	var payloadConfig builder.PayloadConfig

	payloadConfig, err = builder.NewPayloadConfigFromFile(configData)

	if err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}

	//Process payload file
	payloadFilepath, _ := filepath.Abs(Payload)
	payloadData, err := ioutil.ReadFile(payloadFilepath)
	if err != nil {
		logger.Logger.Fatal().Msg(err.Error())
	}

	payloadPath, _, err := payloadConfig.GeneratePayload(payloadData, Executable, true, Parameters)
	if err != nil {
		logger.Logger.Fatal().Msgf("Error during build: %s", err.Error())
	}

	if Output != "" {
		err := os.Rename(payloadPath, Output)
		if err != nil {
			logger.Logger.Fatal().Msgf("Error while moving the payload: %s", err.Error())
		}
		payloadPath = Output
	} 

	logger.Logger.Info().Msgf("Payload can be found here: %s", payloadPath)
}
