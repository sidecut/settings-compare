/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	kv "settings-compare/keyvalue"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff two json settings files, outputting the overrides found in the second file",
	Long: `The first file is considered the base settings file.  The second is
considered the override.  Therefore, the diffs will be the *override*,
i.e. only news settings and changed settings will be output.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			panic("Two file names required")
		}

		// Read files
		map1, err := kv.ReadJsonMapFromFile(args[0])
		if err != nil {
			panic(err)
		}
		map2, err := kv.ReadJsonMapFromFile(args[1])
		if err != nil {
			panic(err)
		}

		kvs1 := kv.GetKeyValues(map1, "")
		kvs2 := kv.GetKeyValues(map2, "")
		log.Println("kvs1:")
		println(kv.GetKeyValuesJson(kvs1))
		log.Println("kvs2:")
		println(kv.GetKeyValuesJson(kvs2))

		// Compare map1 and map2
		println("\nComparing map1 and map2")
		diffs := []kv.KeyValue{}
		for k2, v2 := range map2 {
			if map1[k2] != v2 {
				if map1[k2] == nil {
					// Key not in base file
					fmt.Printf("key:%s\tvalue1:%s\tvalue2:%s\n", k2, map1[k2], v2)
					diffs = append(diffs, kv.KeyValue{k2, v2})
				} else if v2 == nil {
					// 	fmt.Printf("key2:%s\tvalue1:%s\n", k2, map1[k2])
				} else {
					// Also overridden key
					fmt.Printf("key:%s\tvalue1:%s\tvalue2:%s\n", k2, map1[k2], v2)
					diffs = append(diffs, kv.KeyValue{k2, v2})
				}
			}
		}

		println("\nDiffs as a JSON override file:")
		println(kv.GetKeyValuesJson(diffs))

		println("\nNormalized JSON:")
		mNormal, err := kv.MapFromKeyValues(diffs)
		if err != nil {
			panic(err)
		}
		println()
		s, _ := kv.PrettyPrint(mNormal)
		fmt.Printf("%v\n", s)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
