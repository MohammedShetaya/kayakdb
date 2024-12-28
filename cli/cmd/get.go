package cmd

import (
	"github.com/MohammedShetaya/kayakdb/types"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the value of a key",
	Long: `Retrieve the value associated with a key from the database.

This command sends a request to the kayakdb server to fetch the value
associated with the specified key. You need to provide the key name as an
argument. For example:

  kayakctl get myKey

You can also specify the server hostname and port using the global flags:
  -d, --hostname  Set the server hostname (default: "localhost")
  -p, --port      Set the server port (default: "6323")`,
	Args: cobra.ExactArgs(1),
	Run:  commandHandler,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func commandHandler(_ *cobra.Command, args []string) {
	key, _ := types.ConvertStringKeyToDataType(args[0])
	payload := types.Payload{
		Headers: types.Headers{
			Path: "/get",
		},
		Data: []types.KeyValue{
			{Key: key},
		},
	}

	SendRequest(hostname, port, payload, logger)
}
