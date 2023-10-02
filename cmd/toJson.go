/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	kv "settings-compare/keyvalue"

	"github.com/spf13/cobra"
)

// toJsonCmd represents the toJson command
var toJsonCmd = &cobra.Command{
	Use:   "tojson",
	Short: "Convert a file to proper json, i.e. no colons in keys",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("toJson requires exactly 1 argument")
			return
		}

		m, err := kv.ReadJsonMapFromFile(args[0])
		if err != nil {
			panic(err)
		}

		kvs := kv.GetKeyValuesFromMap(m, "")
		// log.Println("kvs1:")
		// println(kv.PrettyPrint(kvs))
		// println()

		mNormal, err := kv.GetMapFromKeyValues(kvs)
		if err != nil {
			panic(err)
		}
		s, _ := kv.PrettyPrint(mNormal)
		fmt.Printf("%v\n", s)
	},
}

func init() {
	rootCmd.AddCommand(toJsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toJsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toJsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
