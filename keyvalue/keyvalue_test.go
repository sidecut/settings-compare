package keyvalue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeToJson_Simple(t *testing.T) {
	kvs := []KeyValue{{"a", "b"}, {"b", "c"}}
	m, err := MapFromKeyValues(kvs)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(m), "Should have 2 keys")
	assert.Equal(t, m["a"], "b")
	assert.Equal(t, m["b"], "c")
}

func TestNormalizeToJson_1Key(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"a:b", "c"}}
	m, err := MapFromKeyValues(kvs)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(m), "Should have 1 key")
}

func TestNormalizeToJson_2Keys(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"a:b", "c"}, {"b:a", "bb"}, {"b:b", "cc"}}
	m, err := MapFromKeyValues(kvs)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(m), "Should have 2 keys")
}

func TestNormalizeToJson_3Keys_1notNested(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"1", "one"}, {"a:b", "c"}, {"b:a", "bb"}, {"b:b", "cc"}}
	m, err := MapFromKeyValues(kvs)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(m), "Should have 3 keys")

	kvs2 := GetKeyValues(m, "")
	assert.Equal(t, 5, len(kvs2))
}

func TestNormalizeToJson_0Keys(t *testing.T) {
	kvs := []KeyValue{}
	m, err := MapFromKeyValues(kvs)

	assert.Nil(t, err)

	assert.Equal(t, 0, len(m), "Should have 0 keys")
}

func TestLoadKeyValuesFromFile(t *testing.T) {
	cases := []struct {
		filename  string
		keycount  int
		keyvalues []KeyValue
	}{
		{"../test-files/base.json", 11,
			[]KeyValue{
				{"this-is-null", nil}, {"this-is-number", 123.0}, {"this-is-boolean", true},
				{"a:a", "aa"}, {"a:b:a", "aba"},
				{"top-level", "top-level-value"},
				{"compound:key", "compound-value"},
				{"b:a", "ba"}, {"b:b:a", "bba"},
				{"delete-me", "base value"}, {"delete-me-2", "base value 2"}}},
		{"../test-files/base.override.json", 7,
			[]KeyValue{
				{"top-level", "overridden-value"},
				{"compound:key", "compound-value"},
				{"b:a", "bleah"}, {"b:c:d", "new-nested-value"}, {"b:c:e", "e"},
				{"delete-me", nil}, {"delete-me-2", ""}}},
	}

	for _, tc := range cases {
		m, err := ReadJsonMapFromFile(tc.filename)

		assert.Nil(t, err)

		kvs := GetKeyValues(m, "")

		assert.Equalf(t, tc.keycount, len(kvs), "Expected %v, got %v", tc.keycount, len(kvs))
		assert.ElementsMatch(t, tc.keyvalues, kvs)
	}
}

func TestPrettyPrint(t *testing.T) {
	kvs := []KeyValue{
		{"this-is-null", nil}, {"this-is-number", 123.0}, {"this-is-boolean", true},
		{"a:a", "aa"}, {"a:b:a", "aba"},
		{"top-level", "top-level-value"},
		// {"empty-map", map[string]interface{}{}},
		{"compound:key", "compound-value"},
		{"b:a", "ba"}, {"b:b:a", "bba"},
		{"delete-me", "base value"}}
	if m, err := MapFromKeyValues(kvs); err != nil {
		panic(err)
	} else {
		// fmt.Printf("%v\n", m)
		if s, err := PrettyPrint(m); err != nil {
			panic(err)
		} else {
			// fmt.Println(s)
			assert.Equal(t, s, `{
  "a": {
    "a": "aa",
    "b": {
      "a": "aba"
    }
  },
  "b": {
    "a": "ba",
    "b": {
      "a": "bba"
    }
  },
  "compound": {
    "key": "compound-value"
  },
  "delete-me": "base value",
  "this-is-boolean": true,
  "this-is-null": null,
  "this-is-number": 123,
  "top-level": "top-level-value"
}`)
		}
	}
}
