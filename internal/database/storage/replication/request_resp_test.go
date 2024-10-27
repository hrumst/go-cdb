package replication

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestResp(t *testing.T) {
	type testObj struct {
		Field1 string
	}

	test1 := testObj{"Test"}

	byteData, err := respReqEncode(&test1)
	assert.NoError(t, err)
	assert.NotEmpty(t, byteData)

	var test2 testObj
	err = respReqDecode(&test2, byteData)
	assert.NoError(t, err)
	assert.Equal(t, test1, test2)
}
