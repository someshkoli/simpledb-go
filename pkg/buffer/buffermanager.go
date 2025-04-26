package buffer

import (
	"time"

	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
	"github.com/someshkoli/simpledb-go/pkg/metrics"
)

type BufferManager struct {
	buffers         []Buffer
	numAvailable    int
	maxTime         time.Duration
	stats           *metrics.InternalStatistics
	freeBufNotifier chan (int)
}

func NewBufferManager(fm *fs.FileManager, lm *log.LogManager, bufPoolSize int, maxTime time.Duration) *BufferManager {
	pool := make([]Buffer, bufPoolSize)
	for n := range bufPoolSize {
		pool[n] = *NewBuffer(fm, lm)
	}
	return &BufferManager{
		buffers:         pool,
		numAvailable:    bufPoolSize,
		maxTime:         maxTime,
		stats:           metrics.NewInternalStatistics("bufferManager"),
		freeBufNotifier: make(chan int),
	}
}

func (fm *BufferManager) Stats() map[string]int {
	return fm.stats.Get()
}

func (bm *BufferManager) Available() int {
	return bm.numAvailable
}

func (bm *BufferManager) FlushAll(txnNum int) {
	for i, b := range bm.buffers {
		if b.ModifyingTxn() == txnNum {
			b.flush()
		}
	}
}

func (bm *BufferManager) findBufferIndex(buff *Buffer) int {
	for n, buff := range bm.buffers {
		if buff.blk.Equals(buff.blk) {
			return n
		}
	}
	return -1
}

func (bm *BufferManager) Unpin(buff *Buffer) {
	buff.UnPin()
	if !buff.IsPinned() {
		bm.numAvailable = bm.numAvailable + 1
		bm.freeBufNotifier <- bm.findBufferIndex(buff)
	}
}

func (bm *BufferManager) waitingTooLong(startTime time.Time) bool {
	return time.Now().Sub(startTime) > time.Duration(bm.maxTime)
}

func (bm *BufferManager) Pin(blk *fs.BlockId) *Buffer {
	buff := bm.tryToPin(blk)
	ticker := time.NewTicker(bm.maxTime)
	if buff != nil {
		return buff
	}
	select {
	case <-ticker.C:
		return nil
	case bno := <-bm.freeBufNotifier:
		b := bm.buffers[bno]
		b.Lock()
		defer b.Unlock()
		if !b.IsPinned() {
			b.Pin()
			b.AssignToBlock(blk)
			return &b
		}
	}
	return nil
}

func (bm *BufferManager) tryToPin(blk *fs.BlockId) *Buffer {

	buff := bm.findExistingBuffer(blk)
	if buff == nil {
		buff = bm.chooseUnpinnedBuffer()
		if buff == nil {
			return nil
		}
		buff.AssignToBlock(blk)
	}
	if buff.IsPinned() {
		bm.numAvailable--
	}
	buff.Pin()
	return buff
}

func (bm *BufferManager) findExistingBuffer(blk *fs.BlockId) *Buffer {
	for _, b := range bm.buffers {
		blk := b.blk
		if blk != nil && blk.Equals(blk) {
			return &b
		}
	}
	return nil
}

func (bm *BufferManager) chooseUnpinnedBuffer() *Buffer {
	for _, b := range bm.buffers {
		if !b.IsPinned() {
			return &b
		}
	}
	return nil
}
