package buffer

import (
	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
)

type Buffer struct {
	filemanager *fs.FileManager
	logManager  *log.LogManager
	contents    *fs.Page
	blk         *fs.BlockId
	pins        int
	txnNumber   int
	lsn         int
}

func NewBuffer(fm *fs.FileManager, lm *log.LogManager, contents *fs.Page, blk *fs.BlockId) *Buffer {
	return &Buffer{
		filemanager: fm,
		logManager:  lm,
		contents:    contents,
		blk:         blk,
		pins:        0,
		txnNumber:   -1,
		lsn:         -1,
	}
}

func (b *Buffer) Contents() *fs.Page {
	return b.contents
}

func (b *Buffer) Block() *fs.BlockId {
	return b.blk
}

func (b *Buffer) SetModified(txnNum, lsn int) {
	b.txnNumber = txnNum
	if lsn >= 0 {
		b.lsn = lsn
	}
}

func (b *Buffer) IsPinned() bool {
	return b.pins > 0
}

func (b *Buffer) ModifyingTxn() int {
	return b.txnNumber
}

func (b *Buffer) AssignToBlock(blk *fs.BlockId) {
	b.flush()
	b.blk = blk
	b.filemanager.Read(blk, b.contents)
	b.pins = 0
}

func (b *Buffer) flush() {
	if b.txnNumber >= 0 {
		b.logManager.Flush(b.lsn)
		b.filemanager.Write(b.blk, b.contents)
		b.txnNumber = -1
	}
}

func (b *Buffer) Pin() {
	b.pins = b.pins + 1
}

func (b *Buffer) UnPin() {
	b.pins = b.pins - 1
}
