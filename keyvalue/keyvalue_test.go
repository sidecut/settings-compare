package keyvalue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeToJson_Simple(t *testing.T) {
	kvs := []KeyValue{{"a", "b"}, {"b", "c"}}
	m, err := MapFromKeyValues(kvs, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(m), "Should have 2 keys")
	assert.Equal(t, m["a"], "b")
	assert.Equal(t, m["b"], "c")
}

func TestNormalizeToJson_1Key(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"a:b", "c"}}
	m, err := MapFromKeyValues(kvs, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(m), "Should have 1 key")
}

func TestNormalizeToJson_2Keys(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"a:b", "c"}, {"b:a", "bb"}, {"b:b", "cc"}}
	m, err := MapFromKeyValues(kvs, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(m), "Should have 2 keys")
}

func TestNormalizeToJson_3Keys_1notNested(t *testing.T) {
	kvs := []KeyValue{{"a:a", "b"}, {"1", "one"}, {"a:b", "c"}, {"b:a", "bb"}, {"b:b", "cc"}}
	m, err := MapFromKeyValues(kvs, nil)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(m), "Should have 3 keys")

	kvs2 := GetKeyValues(m, "")
	assert.Equal(t, 5, len(kvs2))
}

func TestNormalizeToJson_0Keys(t *testing.T) {
	kvs := []KeyValue{}
	m, err := MapFromKeyValues(kvs, nil)

	assert.Nil(t, err)

	assert.Equal(t, 0, len(m), "Should have 0 keys")
}
