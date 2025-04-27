package transaction

import (
	"github.com/someshkoli/simpledb-go/pkg/buffer"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type RecoveryManager struct {
	logManager    *log.LogManager
	bufferManager *buffer.BufferManager
	transaction   *TransactionManager
	txnNumber     int
}

func NewRecoveryManager(
	txn *TransactionManager,
	txnNumber int,
	lm *log.LogManager,
	bm *buffer.BufferManager,
) *RecoveryManager {
	StartLogRecord{}.WriteToLog(lm, txnNumber)
	return &RecoveryManager{
		logManager:    lm,
		bufferManager: bm,
		transaction:   txn,
		txnNumber:     txnNumber,
	}
}

func (rm *RecoveryManager) Commit() {
	rm.bufferManager.FlushAll(rm.txnNumber)
	lsn := CommitLogRecord{}.WriteToLog(rm.logManager, rm.txnNumber)
	rm.logManager.Flush(lsn)
}

func (rm *RecoveryManager) Rollback() {
	// do rollback
	rm.bufferManager.FlushAll(rm.txnNumber)
	lsn := RollbackLogRecord{}.WriteToLog(rm.logManager, rm.txnNumber)
	rm.logManager.Flush(lsn)
}

func (rm *RecoveryManager) Recover(activeTransactions []int) {
	// do recover
	rm.bufferManager.FlushAll(rm.txnNumber)
	lsn := CheckpointLogRecord{}.WriteToLog(rm.logManager, append(activeTransactions, rm.txnNumber))
	rm.logManager.Flush(lsn)
}
