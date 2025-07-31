package cmd

import (
	"github.com/MohammedShetaya/kayakdb/types"
	"github.com/spf13/cobra"
	"os"
)

var (
	hostname string
	port     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kayakctl",
	Short: "A cli to interact with kayakdb server",
	Long:  `A cli to interact with kayakdb server.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&hostname, "hostname", "d", "localhost", "Hostname of the server")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port of the server")
	types.RegisterDataTypes()
}
