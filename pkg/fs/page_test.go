package fs

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetSetInt(t *testing.T) {
	buf := make([]byte, 64)
	page := NewPageFromBytes(buf)
	n := int64(33)
	offset := 10
	page.SetInt(int64(offset), n)
	assert.Equal(t, page.GetInt(int64(offset)), n)
}

func TestGetSetString(t *testing.T) {
	buf := make([]byte, 64)
	page := NewPageFromBytes(buf)

	s := "abcdefgh"
	offset := 20

	page.SetString(int64(offset), s)
	assert.Equal(t, s, page.GetString(int64(offset)))
}
