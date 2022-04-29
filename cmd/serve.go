package cmd

import (
	"github.com/guervild/uru/pkg/api"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a server with an api enpoint /generate.",
	Long:  `Start a server with an api enpoint /generate to generate payload. Listen on 0.0.0.0:8081 by default, can be changed with -a/--addr`,
	Run:   Serve,
}

var Addr string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&Addr, "addr", "a", ":8081", "ip:port to listen on the api server (e.g: 127.0.0.1:3000)")
}

func Serve(cmd *cobra.Command, args []string) {

	a := api.App{}
	a.Initialize()
	a.Run(Addr)
}
