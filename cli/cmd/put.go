package cmd

import (
	"github.com/MohammedShetaya/kayakdb/types"
	"github.com/spf13/cobra"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Set the value of a key",
	Long: `Store a key-value pair in the database.

This command sends a request to the kayakdb server to store the specified
key-value pair. You need to provide both the key name and value as arguments.
For example:

  kayakctl put myKey myValue

You can also specify the server hostname and port using the global flags:
  -d, --hostname  Set the server hostname (default: "localhost")
  -p, --port      Set the server port (default: "6323")`,
	Args: cobra.ExactArgs(2),
	Run:  putCommandHandler,
}

func init() {
	rootCmd.AddCommand(putCmd)
}

func putCommandHandler(_ *cobra.Command, args []string) {
	key, err := ConvertStringToDataType(args[0])
	if err != nil {
		FormatDataTypeError(args[0], err, "key")
	}

	value, err := ConvertStringToDataType(args[1])
	if err != nil {
		FormatDataTypeError(args[1], err, "value")
	}

	payload := types.Payload{
		Headers: types.Headers{
			Path: types.String("/put"),
		},
		Data: []types.Type{
			types.KeyValue{
				Key:   key,
				Value: value,
			},
		},
	}

	SendRequest(hostname, port, payload)
}
