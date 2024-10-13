package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUuidGenerator(t *testing.T) {
	uuid1 := GenerateUUIDv7()
	assert.Len(t, uuid1, 36)

	uuid2 := GenerateUUIDv7()
	assert.NotEqual(t, uuid1, uuid2)
}
