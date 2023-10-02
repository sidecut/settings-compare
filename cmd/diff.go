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

		flatmap1 := kv.GetFlatMap(map1, "")
		flatmap2 := kv.GetFlatMap(map2, "")
		log.Println("kvs1:")
		s, err := kv.PrettyPrint(flatmap1)
		if err != nil {
			panic(err)
		}
		println(s)

		log.Println("kvs2:")
		s, err = kv.PrettyPrint(flatmap2)
		if err != nil {
			panic(err)
		}
		println(s)

		// Compare map1 and map2
		println("\nComparing map1 and map2")
		diffs := []kv.KeyValue{}
		for k2, v2 := range flatmap2 {
			if flatmap1[k2] != v2 {
				if flatmap1[k2] == nil {
					// Key not in base file
					println(fmt.Sprintf("key:%s\tvalue1:%s\tvalue2:%s", k2, flatmap1[k2], v2))
					diffs = append(diffs, kv.KeyValue{Key: k2, Value: v2})
				} else {
					// Overridden key
					println(fmt.Sprintf("key:%s\tvalue1:%s\tvalue2:%s", k2, flatmap1[k2], v2))
					diffs = append(diffs, kv.KeyValue{Key: k2, Value: v2})
				}
			}
		}

		println("\nDiffs as a JSON override file:")
		fm := kv.GetFlatMapFromKeyValues(diffs)
		println(kv.PrettyPrint(fm))

		println("\nNormalized JSON:")
		mNormal, err := kv.GetMapFromKeyValues(diffs)
		if err != nil {
			panic(err)
		}
		s, err = kv.PrettyPrint(mNormal)
		if err != nil {
			panic(err)
		}
		println(s)
		fmt.Println(s)
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
