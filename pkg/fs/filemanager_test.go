package fs

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewFilemanager(t *testing.T) {
	dbDirectory := "./test-db"
	NewFileManager(dbDirectory, 64)

	_, err := os.Stat(dbDirectory)
	assert.NilError(t, err)
}

func TestReadWrite(t *testing.T) {
	dbDirectory := "./test-db"
	fm := NewFileManager(dbDirectory, 128)

	blk := NewBlockId("testfile", 0)
	p1 := NewPageFromBlkSize(fm.blkSize)

	pos1 := 88
	s := "abcdefghijkl"
	p1.SetString(int64(pos1), s)

	size := MaxLength(len(s))
	pos2 := pos1 + size

	n := 345
	p1.SetInt(int64(pos2), int64(n))

	err := fm.Write(blk, p1)
	assert.NilError(t, err)

	p2 := NewPageFromBlkSize(fm.blkSize)
	fm.Read(blk, p2)

	assert.Equal(t, int64(n), p2.GetInt(int64(pos2)))
	assert.Equal(t, s, p2.GetString(int64(pos1)))
}
