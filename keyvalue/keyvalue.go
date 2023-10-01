package keyvalue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type KeyValue struct {
	Key   string
	Value interface{}
}

func ReadJsonMapFromFile(fname string) (map[string]interface{}, error) {
	file1Bytes, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(file1Bytes, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetKeyValues(m map[string]interface{}, prefix string) []KeyValue {
	var keyValues []KeyValue

	for k, v := range m {
		switch v.(type) {
		case string, bool, float64, nil:
			// log.Printf("key:%s\tvalue:%s\n", k, v)
			keyValues = append(keyValues, KeyValue{prefix + k, v})
		default:
			// log.Printf("key:%s\t***value is an object, presumably a map\n", k)
			childKeyValues := GetKeyValues(v.(map[string]interface{}), prefix+k+":")
			keyValues = append(keyValues, childKeyValues...)
		}
	}

	return keyValues
}

func GetFlatMap(m map[string]interface{}, prefix string) map[string]interface{} {
	kvs := GetKeyValues(m, prefix)

	flatmap := make(map[string]interface{})
	for _, kv := range kvs {
		flatmap[kv.Key] = kv.Value
	}

	return flatmap

}

func GetKeyValuesJson(kvs []KeyValue) string {
	// keys := getKeys(kvs)
	// slices.Sort(keys)

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

// func getKeys(kvs []KeyValue) []string {
// 	keys := make([]string, len(kvs))
// 	for i, kv := range kvs {
// 		keys[i] = kv.Key
// 	}
// 	return keys
// }

func MapFromKeyValues(kvs []KeyValue) (map[string]interface{}, error) {
	m := make(map[string]interface{})
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

func PrettyPrint[T any](m T) (string, error) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
