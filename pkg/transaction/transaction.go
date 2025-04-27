package transaction

import (
	"github.com/someshkoli/simpledb-go/pkg/buffer"
	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type TransactionManager struct {
	fileManager   *fs.FileManager
	logManager    *log.LogManager
	bufferManager *buffer.BufferManager
}

func NewTransactionManager(fm *fs.FileManager, lm *log.LogManager, bm *buffer.BufferManager) *TransactionManager {
	return nil
}

func (tm *TransactionManager) Commit() {
	// implement
}

func (tm *TransactionManager) Rollback() {
	// implement
}
