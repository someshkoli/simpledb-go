package log

import (
	"os"
	"strconv"
	"testing"

	"github.com/someshkoli/simpledb-go/pkg/fs"
	"gotest.tools/v3/assert"
)

func TestLogManager(t *testing.T) {
	blkSize := 256
	dbDirectory := "../../test-db"
	fm := fs.NewFileManager(dbDirectory, blkSize)
	logFile := dbDirectory + "/log"
	lm := NewLogManager(fm, logFile)
	createRecords(lm, 1, 35)
	lm.Flush(35)
	iter := lm.Iterator()

	initial := 100 + 35
	for iter.HasNext() {
		rec := iter.Next()
		p := fs.NewPageFromBytes(rec)
		s := p.GetString(0)
		pos := fs.MaxLength(len(s))
		val := p.GetInt(int64(pos))
		assert.Equal(t, val, int64(initial))
		initial = initial - 1
	}
	createRecords(lm, 36, 70)
	lm.Flush(65)
	err := os.Truncate(logFile, 200)
	assert.NilError(t, err)
}

func createRecords(lm *LogManager, start, end int) {
	for i := start; i <= end; i++ {
		rec := createLogRecord("record"+strconv.Itoa(i), i+100)
		lm.Append(rec)
	}
}

func createLogRecord(s string, n int) []byte {
	npos := fs.MaxLength(len(s))
	b := make([]byte, npos+fs.INT64_SIZE)
	p := fs.NewPageFromBytes(b)
	p.SetString(0, s)
	p.SetInt(int64(npos), int64(n))
	return b

}
