package log

import "github.com/someshkoli/simpledb-go/pkg/fs"

type LogIterator struct {
	fileManager *fs.FileManager
	// use buffer page instead of this because same page might already
	// be present in existing buffer.
	blk        *fs.BlockId
	p          *fs.Page
	currentPos int
	boundary   int
}

func moveToBlock(fm *fs.FileManager, blk *fs.BlockId, p *fs.Page) (int, int) {
	fm.Read(blk, p)
	boundary := p.GetInt(0)
	currentPos := boundary
	return int(boundary), int(currentPos)
}

func NewLogIterator(fm *fs.FileManager, blk *fs.BlockId) *LogIterator {
	p := fs.NewPageFromBlkSize(fm.BlkSize)
	boundary, currentPos := moveToBlock(fm, blk, p)
	return &LogIterator{
		fileManager: fm,
		blk:         blk,
		p:           p,
		currentPos:  currentPos,
		boundary:    boundary,
	}
}

func (i *LogIterator) HasNext() bool {
	return i.currentPos < i.fileManager.BlkSize || i.blk.Number() > 0
}

func (i *LogIterator) Next() []byte {
	if i.currentPos == i.fileManager.BlkSize {
		i.p = fs.NewPageFromBlkSize(i.fileManager.BlkSize)
		i.blk = fs.NewBlockId(i.blk.FileName(), i.blk.Number()-1)
		i.boundary, i.currentPos = moveToBlock(i.fileManager, i.blk, i.p)
	}
	rec := i.p.GetBytes(int64(i.currentPos))
	i.currentPos = i.currentPos + fs.INT64_SIZE + len(rec)
	return rec
}
