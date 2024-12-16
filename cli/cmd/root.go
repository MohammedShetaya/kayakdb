package cmd

import (
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	hostname string
	port     string
	logger   *zap.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kayakctl",
	Short: "A cli to interact with kayakdb server",
	Long:  `A cli to interact with kayakdb server.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize the logger once before any command runs
		if logger == nil {
			logger = log.InitLogger()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Ensure the logger is flushed after execution
		if logger != nil {
			_ = logger.Sync()
		}
	},
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
	api.InitProtocol()
}
