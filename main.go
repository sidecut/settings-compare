package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

var (
	file1 = pflag.StringP("base", "b", "", "base filename, e.g. appsettings.json")
	file2 = pflag.StringP("override", "o", "", "override filename, e.g. appsettings.development.json")
)

type KeyValue struct {
	Key   string
	Value interface{}
}

func main() {
	pflag.Parse()

	if *file1 == "" {
		panic("base file must be specified")
	}

	if *file2 == "" {
		panic("override file must be specified")
	}

	map1 := make(map[string]interface{})
	map2 := make(map[string]interface{})

	// Read file1 into map1
	file1Bytes, err := os.ReadFile(*file1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file1Bytes, &map1)
	if err != nil {
		panic(err)
	}

	// Read file2 into map2
	file2Bytes, err := os.ReadFile(*file2)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file2Bytes, &map2)
	if err != nil {
		panic(err)
	}

	kvs1 := getKeyValues(map1, "")
	kvs2 := getKeyValues(map2, "")
	log.Println("kvs1:")
	println(getKeyValuesJson(kvs1))
	log.Println("kvs2:")
	println(getKeyValuesJson(kvs2))

	// Compare map1 and map2
	println("\nComparing map1 and map2")
	diffs := []KeyValue{}
	for k2, v2 := range map2 {
		if map1[k2] != v2 {
			if map1[k2] == nil {
				// Key not in base file
				fmt.Printf("key:%s\tvalue1:%s\tvalue2:%s\n", k2, map1[k2], v2)
				diffs = append(diffs, KeyValue{k2, v2})
			} else if v2 == nil {
				// 	fmt.Printf("key2:%s\tvalue1:%s\n", k2, map1[k2])
			} else {
				// Also overridden key
				fmt.Printf("key:%s\tvalue1:%s\tvalue2:%s\n", k2, map1[k2], v2)
				diffs = append(diffs, KeyValue{k2, v2})
			}
		}
	}

	println("\nDiffs as a JSON override file:")
	println(getKeyValuesJson(diffs))

	println("\nNormalized JSON:")
	mNormal, err := mapFromKeyValues(diffs, nil)
	if err != nil {
		panic(err)
	}
	println()
	s, _ := prettyPrint(mNormal)
	fmt.Printf("%v\n", s)
}

func getKeyValues(m map[string]interface{}, prefix string) []KeyValue {
	keyValues := []KeyValue{}

	for k, v := range m {
		switch v.(type) {
		case string:
			// log.Printf("key:%s\tvalue:%s\n", k, v)
			keyValues = append(keyValues, KeyValue{k, v})
		default:
			// log.Printf("key:%s\t***value is an object, presumably a map\n", k)
			childKeyValues := getKeyValues(v.(map[string]interface{}), k+":")
			keyValues = append(keyValues, childKeyValues...)
		}
	}

	return keyValues
}

func getKeyValuesJson(kvs []KeyValue) string {
	keys := getKeys(kvs)
	slices.Sort(keys)

	m := make(map[string]interface{})
	for _, kv := range kvs {
		m[kv.Key] = kv.Value
	}

	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling json: %v", err)
	}

	return string(jsonBytes)
}

func getKeys(kvs []KeyValue) []string {
	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}
	return keys
}

func mapFromKeyValues(kvs []KeyValue, m map[string]interface{}) (map[string]interface{}, error) {
	if m == nil {
		m = make(map[string]interface{})
	}
	for _, kv := range kvs {
		if err := putIntoMap(kv, m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func putIntoMap(kv KeyValue, m map[string]interface{}) error {
	// This function drills down into a map non-recursively, i.e. `m` will be assigned new values throughout
	keyParts := strings.Split(kv.Key, ":")
	for i, keyPart := range keyParts {
		// println("i, keyPart", i, keyPart)
		if i+1 == len(keyParts) {
			// leaf
			// println("Leaf")
			if m[keyPart] != nil {
				return fmt.Errorf("WARN: %v already defined as %v, now getting %v", kv.Key, m[keyPart], kv.Value)
			} else {
				m[keyPart] = kv.Value
			}
		} else {
			// interior
			// println("Interior")
			if m[keyPart] == nil {
				m[keyPart] = make(map[string]interface{})
				// putIntoMap()
				m = m[keyPart].(map[string]interface{})
			} else {
				// drill down
				m = m[keyPart].(map[string]interface{})
			}
		}
	}

	return nil
}

func prettyPrint(m map[string]interface{}) (string, error) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
