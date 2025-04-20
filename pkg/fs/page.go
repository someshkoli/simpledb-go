package fs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

type BlockId struct {
	filename string
	blknum   int
}

func NewBlockId(fileName string, blknum int) *BlockId {
	return &BlockId{
		filename: fileName,
		blknum:   blknum,
	}
}

func (b *BlockId) Equals(in *BlockId) bool {
	return b.filename == in.filename && b.blknum == in.blknum
}

func (b *BlockId) ToString() string {
	return fmt.Sprintf("[file: %s, blocknum: %d]", b.filename, b.blknum)
}

func (b *BlockId) HashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(b.ToString()))
	return h.Sum32()
}
func (b *BlockId) Number() int {
	return b.blknum
}

func (b *BlockId) FileName() string {
	return b.filename
}

type Page struct {
	bb bytes.Buffer
}

func NewPageFromBytes(buf []byte) *Page {
	return &Page{
		bb: *bytes.NewBuffer(buf),
	}
}

func NewPageFromBlkSize(blkSize int) *Page {
	return &Page{
		bb: *bytes.NewBuffer(make([]byte, blkSize)),
	}
}

// hadnle this hardcoding
const INT64_SIZE = 8

func (p *Page) buffer() []byte {
	return p.bb.Bytes()
}

func (p *Page) GetInt(offset int64) int64 {
	return int64(
		binary.BigEndian.Uint64(
			p.buffer()[offset : offset+INT64_SIZE],
		),
	)
}

func (p *Page) SetInt(offset int64, n int64) {
	buf := make([]byte, INT64_SIZE)
	binary.BigEndian.PutUint64(buf, uint64(n))
	copy(p.buffer()[offset:offset+INT64_SIZE], buf)
}

func (p *Page) GetBytes(offset int64) []byte {
	n := p.GetInt(offset)
	offset = offset + INT64_SIZE
	return p.buffer()[offset : offset+n]
}

func (p *Page) SetBytes(offset int64, b []byte) {
	p.SetInt(offset, int64(len(b)))
	offset = offset + INT64_SIZE
	copy(p.buffer()[offset:offset+int64(len(b))], b)
}

func (p *Page) GetString(offset int64) string {
	return string(p.GetBytes(offset))
}

func (p *Page) SetString(offset int64, s string) {
	p.SetBytes(offset, []byte(s))
}

func MaxLength(strlen int) int {
	// implmennt character set here
	return INT64_SIZE + strlen
}
