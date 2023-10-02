/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	kv "settings-compare/keyvalue"

	"github.com/spf13/cobra"
)

// toKeysCmd represents the toKeys command
var toKeysCmd = &cobra.Command{
	Use:   "tokeys",
	Short: "Convert a file to .NET user-secrets form, i.e. colons in keys",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("toKeys requires exactly 1 argument")
			return
		}

		m, err := kv.ReadJsonMapFromFile(args[0])
		if err != nil {
			panic(err)
		}

		fm := kv.GetFlatMap(m, "")
		s, err := kv.PrettyPrint(fm)
		if err != nil {
			panic(err)
		}
		cmd.Printf("%v\n", s)
	},
}

func init() {
	rootCmd.AddCommand(toKeysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toKeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toKeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
