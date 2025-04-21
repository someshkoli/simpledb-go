package buffer

import (
	"time"

	"github.com/someshkoli/simpledb-go/pkg/fs"
	"github.com/someshkoli/simpledb-go/pkg/log"
	"github.com/someshkoli/simpledb-go/pkg/metrics"
	"github.com/someshkoli/simpledb-go/pkg/utils/observer"
)

type BufferManager struct {
	buffer       []Buffer
	numAvailable int
	maxTime      int
	observers    []observer.Observer
	stats        *metrics.InternalStatistics
}

func NewBufferManager(fm *fs.FileManager, lm *log.LogManager, bufPoolSize int, maxTime int) *BufferManager {
	pool := make([]Buffer, bufPoolSize)
	for n := range bufPoolSize {
		pool[n] = *NewBuffer(fm, lm)
	}
	return &BufferManager{
		buffer:       pool,
		numAvailable: bufPoolSize,
		maxTime:      maxTime,
		stats:        metrics.NewInternalStatistics("bufferManager"),
	}
}

func (fm *BufferManager) Stats() map[string]int {
	return fm.stats.Get()
}

func (bm *BufferManager) RegisterObserver(ob observer.Observer) {
	bm.observers = append(bm.observers, ob)
}

func (bm *BufferManager) DeregisterObserver(ob observer.Observer) {
	for i, o := range bm.observers {
		if ob.GetId() == o.GetId() {
			bm.observers[i] = bm.observers[len(bm.observers)-1]
			bm.observers = bm.observers[:len(bm.observers)-1]
		}
	}
}

func (bm *BufferManager) notifyAllObserver() {
	for _, o := range bm.observers {
		o.Update()
	}
}

func (bm *BufferManager) Available() int {
	return bm.numAvailable
}

func (bm *BufferManager) FlushAll(txnNum int) {
	for i, b := range bm.buffer {
		if b.ModifyingTxn() == txnNum {
			b.flush()
		}
	}
}

func (bm *BufferManager) Unpin(buff *Buffer) {
	buff.UnPin()
	if !buff.IsPinned() {
		bm.numAvailable = bm.numAvailable + 1
		bm.notifyAllObserver()
	}
}

func (bm *BufferManager) waitingTooLong(startTime time.Time) bool {
	return time.Now().Sub(startTime) > time.Duration(bm.maxTime)
}

func (bm *BufferManager) Pin(blk *fs.BlockId) *Buffer {
	now := time.Now()
	buff := bm.tryToPin(blk)
	for buff != nil && !bm.waitingTooLong(now) {
		// (bm.maxTime) // accept notification
		// use channels instead of observers
	}

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
	for _, b := range bm.buffer {
		blk := b.blk
		if blk != nil && blk.Equals(blk) {
			return &b
		}
	}
	return nil
}

func (bm *BufferManager) chooseUnpinnedBuffer() *Buffer {
	for _, b := range bm.buffer {
		if !b.IsPinned() {
			return &b
		}
	}
	return nil
}
