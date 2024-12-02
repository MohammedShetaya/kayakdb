/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the value of a key",
	Long: `Get the value of a key by providing the kay name:
get <kay_name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := ConvertToDataType(args[0])
		payload := api.Payload{
			// TODO: add other headers dynamically in the cli
			Headers: api.Headers{
				Path: "/get",
			},
			Data: []api.KeyValue{
				{Key: key},
			},
		}
		fmt.Println(payload.String())
		SendRequest(hostname, port, payload)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
