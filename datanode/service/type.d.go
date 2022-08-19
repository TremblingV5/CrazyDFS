package DNService

import (
	"bufio"
	"os"
)

type Block struct {
	Name      string
	Offset    int64
	ChunkSize int
	Reader    *bufio.Reader
	Buffer    *[]byte
	Cursor    int
	File      *os.File
	BlockSize int64
}
