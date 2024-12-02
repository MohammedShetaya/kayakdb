package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.Flags().StringVarP(&hostname, "hostname", "d", "localhost", "Hostname of the server")
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port of the server")
}
