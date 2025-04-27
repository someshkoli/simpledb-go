package db

import (
	"time"

	"github.com/someshkoli/simpledb-go/pkg/buffer"
	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type Database struct {
	FileManager   *fs.FileManager
	LogManager    *log.LogManager
	BufferManager *buffer.BufferManager
}

func NewDatabase(databaseName string, blkSize int, bufferPoolSize int, bufferPoolMaxTime time.Duration) *Database {

	fm := fs.NewFileManager(databaseName, blkSize)
	lm := log.NewLogManager(fm, databaseName+"_log")
	return &Database{
		FileManager:   fm,
		LogManager:    lm,
		BufferManager: buffer.NewBufferManager(fm, lm, bufferPoolSize, bufferPoolMaxTime),
	}
}
