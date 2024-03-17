package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPtr(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, *(boolPtr(true)))
	assert.Equal(false, *(boolPtr(false)))
}

func TestInt32Ptr(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(int32(10), *(int32Ptr(10)))
}
