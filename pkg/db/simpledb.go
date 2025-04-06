package db

import "github.com/someshkoli/simpledb-go/pkg/fs"

type Database struct {
	FileManager *fs.FileManager
}

func NewDatabase(databaseName string, blkSize int, bufferPoolSize int) *Database {
	return &Database{
		FileManager: fs.NewFileManager(databaseName, blkSize),
	}
}
