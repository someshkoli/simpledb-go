package db

import (
	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type Database struct {
	FileManager *fs.FileManager
	LogManager  *log.LogManager
}

func NewDatabase(databaseName string, blkSize int, bufferPoolSize int) *Database {

	fm := fs.NewFileManager(databaseName, blkSize)
	return &Database{
		FileManager: fm,
		LogManager:  log.NewLogManager(fm, databaseName+"_log"),
	}
}
