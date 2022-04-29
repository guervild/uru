package cmd

import (
	"log"

	"github.com/guervild/uru/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	JsonLog bool
	Debug bool
)

var rootCmd = &cobra.Command{
	Use:   "uru",
	Short: "Payload generation tool for windows.",
	Long:  `Payload generation tool that uses a config file to defined an execution workflow. It helps you to obfuscate and execute a shellcode/executable during your engagement.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.InitLogger(JsonLog, Debug)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&JsonLog, "jsonlog", "", false, "Print logs using json output")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "Print addtionnal debug log")

}
