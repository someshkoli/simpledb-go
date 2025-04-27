package transaction

import (
	"fmt"

	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type LogType int

const CHECKPOINT LogType = 0
const START LogType = 1
const COMMIT LogType = 2
const ROLLBACK LogType = 3
const SETINT LogType = 4
const SETSTRING LogType = 5

type LogRecord interface {
	Op() LogType
	TxnNumber() int
	Undo(txnNumber int)
	ToString() string
}

func CreateLogRecord(bytes []byte) LogRecord {
	p := fs.NewPageFromBytes(bytes)
	switch p.GetInt(0) {
	case int64(CHECKPOINT):
	case int64(START):
	case int64(COMMIT):
		lr := NewCommitLogRecord(p)
		return &lr
	case int64(ROLLBACK):
	case int64(SETINT):
	case int64(SETSTRING):
	}
	return nil
}

type CommitLogRecord struct {
	txnNumber int
}

func NewCommitLogRecord(p *fs.Page) CommitLogRecord {
	tpos := fs.INT64_SIZE
	txnNumber := p.GetInt(int64(tpos))
	return CommitLogRecord{
		txnNumber: int(txnNumber),
	}
}

func (c *CommitLogRecord) Op() LogType {
	return CHECKPOINT
}

func (c *CommitLogRecord) TxnNumber() int {
	return c.txnNumber
}

func (c *CommitLogRecord) Undo(txnNumber int) {
	// implement this properly
	return
}

func (c *CommitLogRecord) ToString() string {
	return fmt.Sprintf("<COMMIT %d>", c.txnNumber)
}

func (c CommitLogRecord) WriteToLog(lm *log.LogManager, txnNumber int) int {
	tpos := fs.INT64_SIZE
	recLen := fs.MaxLength(tpos)
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)
	p.SetInt(0, int64(CHECKPOINT))
	p.SetInt(int64(tpos), int64(txnNumber))
	return lm.Append(rec)
}

type SetStringLogRecord struct {
	txnNumber int
	val       []byte
	oldVal    []byte
	blk       *fs.BlockId
	offset    int
}

func NewSetStringLogRecord(p *fs.Page) SetStringLogRecord {
	tpos := fs.INT64_SIZE
	txnNumber := p.GetInt(int64(tpos))
	fpos := tpos + fs.INT64_SIZE
	filename := p.GetString(int64(fpos))
	bpos := fpos + fs.MaxLength(len(filename))
	blkNumber := p.GetInt(int64(bpos))
	blk := fs.NewBlockId(filename, int(blkNumber))
	opos := bpos + fs.INT64_SIZE
	offset := p.GetInt(int64(opos))
	vpos := opos + fs.INT64_SIZE
	val := p.GetBytes(int64(vpos))
	ovpos := vpos + fs.MaxLength(len(val))
	oldVal := p.GetBytes(int64(ovpos))
	return SetStringLogRecord{
		txnNumber: int(txnNumber),
		val:       val,
		oldVal:    oldVal,
		blk:       blk,
		offset:    int(offset),
	}
}

func (lr *SetStringLogRecord) Op() int {
	return int(SETSTRING)
}

func (lr *SetStringLogRecord) TxnNumber() int {
	return lr.txnNumber
}

func (lr *SetStringLogRecord) ToString() string {
	return fmt.Sprintf(
		"<SETSTRING %d %d %d %s %s>",
		lr.txnNumber,
		lr.blk.Number(),
		lr.offset,
		fs.NewPageFromBytes(lr.oldVal).GetString(int64(lr.offset)),
		fs.NewPageFromBytes(lr.val).GetString(int64(lr.offset)),
	)
}

// needs transaction to be implemented
func (lr *SetStringLogRecord) Undo() {

}

func (lr *SetStringLogRecord) WriteToLog(lm *log.LogManager, txnNum int, blk *fs.BlockId, offset int, val []byte, oldVal []byte) int {
	tpos := fs.INT64_SIZE
	fpos := tpos + fs.INT64_SIZE
	bpos := fpos + fs.MaxLength(len(blk.FileName()))
	opos := bpos + fs.INT64_SIZE
	vpos := opos + fs.INT64_SIZE
	ovpos := vpos + fs.MaxLength(len(val))
	recLen := vpos + fs.MaxLength(len(oldVal))
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)
	p.SetInt(0, int64(SETSTRING))
	p.SetInt(int64(tpos), int64(txnNum))
	p.SetString(int64(fpos), blk.FileName())
	p.SetInt(int64(bpos), int64(blk.Number()))
	p.SetInt(int64(opos), int64(offset))
	p.SetBytes(int64(vpos), val)
	p.SetBytes(int64(ovpos), oldVal)
	return lm.Append(rec)
}

type StartLogRecord struct {
	txnNumber int
}

// Constructor for StartLogRecord
func NewStartLogRecord(p *fs.Page) StartLogRecord {
	tpos := fs.INT64_SIZE
	txnNumber := p.GetInt(int64(tpos))
	return StartLogRecord{
		txnNumber: int(txnNumber),
	}
}

func (lr *StartLogRecord) Op() LogType {
	return START
}

func (lr *StartLogRecord) TxnNumber() int {
	return lr.txnNumber
}

func (lr *StartLogRecord) Undo(txnNumber int) {
	// START log records typically do not require undo logic
	return
}

func (lr *StartLogRecord) ToString() string {
	return fmt.Sprintf("<START %d>", lr.txnNumber)
}

// Method to write a StartLogRecord to the log
func (lr StartLogRecord) WriteToLog(lm *log.LogManager, txnNumber int) int {
	tpos := fs.INT64_SIZE
	recLen := tpos + fs.INT64_SIZE
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)
	p.SetInt(0, int64(START)) // Log type
	p.SetInt(int64(tpos), int64(txnNumber))
	return lm.Append(rec)
}

type SetIntLogRecord struct {
	txnNumber int
	blk       *fs.BlockId
	offset    int
	val       []byte
	oldVal    []byte
}

// Constructor for SetIntLogRecord
func NewSetIntLogRecord(p *fs.Page) SetIntLogRecord {
	tpos := fs.INT64_SIZE
	txnNumber := p.GetInt(int64(tpos))
	fpos := tpos + fs.INT64_SIZE
	filename := p.GetString(int64(fpos))
	bpos := fpos + fs.MaxLength(len(filename))
	blkNumber := p.GetInt(int64(bpos))
	blk := fs.NewBlockId(filename, int(blkNumber))
	opos := bpos + fs.INT64_SIZE
	offset := p.GetInt(int64(opos))
	vpos := opos + fs.INT64_SIZE
	val := p.GetBytes(int64(vpos))
	ovPos := vpos + fs.MaxLength(len(val))
	oldVal := p.GetBytes(int64(ovPos))
	return SetIntLogRecord{
		txnNumber: int(txnNumber),
		blk:       blk,
		offset:    int(offset),
		val:       val,
		oldVal:    oldVal,
	}
}

func (lr *SetIntLogRecord) Op() LogType {
	return SETINT
}

func (lr *SetIntLogRecord) TxnNumber() int {
	return lr.txnNumber
}

func (lr *SetIntLogRecord) Undo(txnNumber int) {
	// Undo logic for SetIntLogRecord
	// Typically, this would involve restoring the previous value in the block
	// Needs transaction logic to be implemented
}

func (lr *SetIntLogRecord) ToString() string {
	return fmt.Sprintf(
		"<SETINT %d %d %d %d %d>",
		lr.txnNumber,
		lr.blk.Number(),
		lr.offset,
		fs.NewPageFromBytes(lr.oldVal).GetInt(int64(lr.offset)),
		fs.NewPageFromBytes(lr.val).GetInt(int64(lr.offset)),
	)
}

// Method to write a SetIntLogRecord to the log
func (lr *SetIntLogRecord) WriteToLog(lm *log.LogManager, txnNum int, blk *fs.BlockId, offset int, val []byte, oldVal []byte) int {
	tpos := fs.INT64_SIZE
	fpos := tpos + fs.INT64_SIZE
	bpos := fpos + fs.MaxLength(len(blk.FileName()))
	opos := bpos + fs.INT64_SIZE
	vpos := opos + fs.INT64_SIZE
	ovpos := opos + fs.MaxLength(len(val))
	recLen := vpos + fs.MaxLength(len(oldVal))
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)
	p.SetInt(0, int64(SETINT))
	p.SetInt(int64(tpos), int64(txnNum))
	p.SetString(int64(fpos), blk.FileName())
	p.SetInt(int64(bpos), int64(blk.Number()))
	p.SetInt(int64(opos), int64(offset))
	p.SetBytes(int64(vpos), val)
	p.SetBytes(int64(ovpos), oldVal)
	return lm.Append(rec)
}

type RollbackLogRecord struct {
	txnNumber int
}

func NewRollbackLogRecord(p *fs.Page) RollbackLogRecord {
	tpos := fs.INT64_SIZE
	txnNumber := p.GetInt(int64(tpos))
	return RollbackLogRecord{
		txnNumber: int(txnNumber),
	}
}

func (r *RollbackLogRecord) Op() LogType {
	return ROLLBACK
}

func (r *RollbackLogRecord) TxnNumber() int {
	return r.txnNumber
}

func (r *RollbackLogRecord) Undo(txnNumber int) {
	// ROLLBACK log records typically do not require undo logic
	return
}

func (r *RollbackLogRecord) ToString() string {
	return fmt.Sprintf("<ROLLBACK %d>", r.txnNumber)
}

func (r RollbackLogRecord) WriteToLog(lm *log.LogManager, txnNumber int) int {
	tpos := fs.INT64_SIZE
	recLen := fs.MaxLength(tpos)
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)
	p.SetInt(0, int64(ROLLBACK))
	p.SetInt(int64(tpos), int64(txnNumber))
	return lm.Append(rec)
}

// this is a type of Non quiescent checkpoint which avoid database blocking
// while the recovery manager waits for existing transactions to complete.
type CheckpointLogRecord struct {
	activeTxns []int
}

// Constructor for CheckpointLogRecord
func NewCheckpointLogRecord(p *fs.Page) CheckpointLogRecord {
	tpos := fs.INT64_SIZE
	numTxns := p.GetInt(int64(tpos)) // Number of active transactions
	tpos += fs.INT64_SIZE

	activeTxns := make([]int, numTxns)
	for i := 0; i < int(numTxns); i++ {
		activeTxns[i] = int(p.GetInt(int64(tpos)))
		tpos += fs.INT64_SIZE
	}

	return CheckpointLogRecord{
		activeTxns: activeTxns,
	}
}

func (c *CheckpointLogRecord) Op() LogType {
	return CHECKPOINT
}

func (c *CheckpointLogRecord) TxnNumber() int {
	// Checkpoint does not have a single transaction number
	return -1
}

func (c *CheckpointLogRecord) Undo(txnNumber int) {
	// Checkpoint log records do not require undo logic
	return
}

func (c *CheckpointLogRecord) ToString() string {
	return fmt.Sprintf("<CHECKPOINT %v>", c.activeTxns)
}

func (c CheckpointLogRecord) WriteToLog(lm *log.LogManager, activeTxns []int) int {
	tpos := fs.INT64_SIZE
	recLen := tpos + fs.INT64_SIZE + len(activeTxns)*fs.INT64_SIZE
	rec := make([]byte, recLen)
	p := fs.NewPageFromBytes(rec)

	p.SetInt(0, int64(CHECKPOINT)) // Log type
	p.SetInt(int64(tpos), int64(len(activeTxns)))
	tpos += fs.INT64_SIZE

	for _, txn := range activeTxns {
		p.SetInt(int64(tpos), int64(txn))
		tpos += fs.INT64_SIZE
	}

	return lm.Append(rec)
}
