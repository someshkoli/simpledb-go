package fs

import (
	"fmt"
	"os"
	"strings"

	"github.com/someshkoli/simpledb-go/pkg/metrics"
)

type STATISTICS_KEYS string

const FILE_STAT_READ = "file_stat_read"
const BLOCK_READ = "block_read"
const BLOCK_WRITE = "block_write"

type FileManager struct {
	dbDirectory string
	BlkSize     int
	IsNew       bool
	stats       *metrics.InternalStatistics
}

func NewFileManager(dbDirectory string, blkSize int) *FileManager {
	isNew := false
	if _, err := os.Stat(dbDirectory); err != nil {
		os.Mkdir(dbDirectory, 0777)
		isNew = true
	}

	f, err := os.ReadDir(dbDirectory)
	if err != nil {
		panic(fmt.Errorf("unable to read db directory, %w", err))
	}

	for _, n := range f {
		if strings.HasPrefix(n.Name(), "temp_") {
			os.Remove(n.Name())
		}
	}

	return &FileManager{
		dbDirectory: dbDirectory,
		BlkSize:     int(blkSize),
		IsNew:       isNew,
		stats:       metrics.NewInternalStatistics("filemanager"),
	}
}

func (fm *FileManager) Read(blk *BlockId, p *Page) error {
	f, err := os.OpenFile(blk.filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	offset := blk.blknum * fm.BlkSize
	f.Seek(int64(offset), 0)
	f.Read(p.buffer())
	fm.stats.Increment(BLOCK_READ, 1)

	return nil
}

func (fm *FileManager) Write(blk *BlockId, p *Page) error {
	f, err := os.OpenFile(blk.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	offset := blk.blknum * fm.BlkSize
	_, err = f.WriteAt(p.buffer(), int64(offset))
	if err != nil {
		return err
	}
	fm.stats.Increment(BLOCK_WRITE, 1)

	return nil
}

func (fm *FileManager) Append(filename string) (*BlockId, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	newBlkNum := f.Size() / int64(fm.BlkSize)
	blk := NewBlockId(filename, int(newBlkNum))

	nf, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer nf.Close()

	nf.Seek(f.Size(), 0)
	nf.Write(make([]byte, fm.BlkSize))
	fm.stats.Increment(BLOCK_WRITE, 1)

	return blk, nil
}

func (fm *FileManager) Length(filename string) (int, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return -1, err
	}
	fm.stats.Increment(FILE_STAT_READ, 1)

	return int(f.Size()) / fm.BlkSize, nil
}

func (fm *FileManager) Stats() map[string]int {
	fm.stats.Increment(FILE_STAT_READ, 1)

	return fm.stats.Get()
}
