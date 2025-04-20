package log

import (
	"os"

	"github.com/someshkoli/simpledb-go/pkg/fs"
)

type LogManager struct {
	fileManager  *fs.FileManager
	logFileName  string
	logPage      *fs.Page
	currentBlk   *fs.BlockId
	latestLSN    int
	lastSavedLSN int
}

func appendNewBlock(fileManager *fs.FileManager, logFile string, logPage *fs.Page) *fs.BlockId {
	blk, err := fileManager.Append(logFile)
	if err != nil {
		panic(err) // unable to open file or some file related error
	}
	logPage.SetInt(0, int64(fileManager.BlkSize)) // write the curent boundary, initially it will be at the end of the page
	err = fileManager.Write(blk, logPage)
	if err != nil {
		panic(err) // file write related error, should panic because log is mandatory component
	}
	return blk
}

func NewLogManager(fileManager *fs.FileManager, logFile string) *LogManager {
	logPage := fs.NewPageFromBlkSize(fileManager.BlkSize)

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err) // means file does not exist or is unavailable
	}
	defer f.Close()

	logSize, err := fileManager.Length(logFile)
	if err != nil {
		panic(err)
	}

	var currentBlk *fs.BlockId
	if logSize == 0 {
		currentBlk = appendNewBlock(fileManager, logFile, logPage)
	} else {
		currentBlk = fs.NewBlockId(logFile, logSize-1)
		fileManager.Read(currentBlk, logPage)
	}
	return &LogManager{
		fileManager:  fileManager,
		logFileName:  logFile,
		currentBlk:   currentBlk,
		logPage:      logPage,
		latestLSN:    0,
		lastSavedLSN: 0,
	}
}

func (lm *LogManager) Flush(lsn int) {
	if lsn >= lm.lastSavedLSN {
		lm.flush()
	}
}

func (lm *LogManager) flush() {
	lm.fileManager.Write(lm.currentBlk, lm.logPage)
	lm.lastSavedLSN = lm.latestLSN
}

func (lm *LogManager) Append(logrec []byte) int {
	boundary := lm.logPage.GetInt(0)
	recSize := len(logrec)
	bytesNeeded := recSize + fs.INT64_SIZE // + int size to store length of the bytes
	if boundary-int64(bytesNeeded) < fs.INT64_SIZE {
		lm.flush()
		logPage := fs.NewPageFromBlkSize(lm.fileManager.BlkSize)
		lm.logPage = logPage
		lm.currentBlk = appendNewBlock(lm.fileManager, lm.logFileName, logPage)
		boundary = lm.logPage.GetInt(0)
	}
	recpos := boundary - int64(bytesNeeded)
	lm.logPage.SetBytes(recpos, logrec)
	lm.logPage.SetInt(0, recpos) // set new boundary at first position
	lm.latestLSN = lm.latestLSN + 1
	return lm.latestLSN
}

func (lm *LogManager) Iterator() *LogIterator {
	return NewLogIterator(lm.fileManager, lm.currentBlk)
}
