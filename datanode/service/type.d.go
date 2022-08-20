package DNService

import (
	"bufio"
	"os"
)

type Block struct {
	ID        string
	Offset    int64
	ChunkSize int
	Reader    *bufio.Reader
	Buffer    *[]byte
	Cursor    int
	File      *os.File
	BlockSize int64
}

type BlockID int64
