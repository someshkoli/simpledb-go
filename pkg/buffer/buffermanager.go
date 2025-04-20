package buffer

type BufferManager struct {
	buffer       []Buffer
	numAvailable int
	maxTime      int
}
