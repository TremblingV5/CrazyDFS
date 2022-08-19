package DNService

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
)

var config, _ = utils.InitNodeConfig(items.DN{}, values.DataNodeConfigPath)

func (b *Block) initBlock(name string, mode string) {
	var err error
	var reader *bufio.Reader
	var file *os.File

	if mode == "r" {
		file, err = os.Open(name)
		reader = bufio.NewReader(file)
	} else if mode == "w" {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	}

	if err != nil {
		utils.WriteLog(
			"error", "can't open file"+err.Error(),
		)
	}

	b.File = file
	b.Name = name
	b.Reader = reader
	b.ChunkSize = int(config.IOSize)
	b.BlockSize = config.BlockSize
	buffer := make([]byte, b.ChunkSize)
	b.Buffer = &buffer
	b.Cursor = -1
	b.Offset = 0
}

func GetBlock(name string, mode string) *Block {
	block := Block{}
	block.initBlock(
		name, mode,
	)
	return &block
}

func (b *Block) HasNextChunk() bool {
	if b.Cursor != -1 {
		return true
	}

	n, err := b.Reader.Read(*b.Buffer)
	if err == io.EOF {
		b.Cursor = -1
		b.File.Close()
		return false
	}
	if err != nil {
		utils.WriteLog(
			"error", "Read file defeat",
			"message", err.Error(),
		)
	}

	b.Cursor = n
	return true
}

func (b *Block) GetNextBlock() (*[]byte, int, error) {
	if b.Cursor == -1 {
		return nil, 0, nil
	}

	n := b.Cursor
	b.Cursor = -1
	return b.Buffer, n, nil
}

func (b *Block) WriteChunk(chunk []byte) error {
	info, err := b.File.Stat()
	if err != nil {
		utils.WriteLog(
			"error", "Get file stat defeat",
			"message", err.Error(),
		)
	}

	currBlockSize := info.Size()
	if b.BlockSize >= int64(len(chunk)+int(currBlockSize)) {
		if _, err := b.File.Write(chunk); err != nil {
			utils.WriteLog(
				"error", "Write chunk defeat",
				"message", err.Error(),
			)
		}
		return nil
	}

	return errors.New("file larger than block size")
}

func (b *Block) Close() error {
	return b.File.Close()
}

func (b *Block) GetFileSize() int64 {
	info, _ := b.File.Stat()
	return info.Size()
}
