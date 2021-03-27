package objects_test

import (
	"github.com/YReshetko/rash-lang/objects"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	key1 := &objects.String{Value: "key"}
	key2 := &objects.String{Value: "key"}
	value1 := &objects.String{Value: "value"}
	value2 := &objects.String{Value: "value"}

	assert.Equal(t, key1.HashKey(), key2.HashKey())
	assert.Equal(t, value1.HashKey(), value2.HashKey())
	assert.NotEqual(t, key1.HashKey(), value1.HashKey())
}
