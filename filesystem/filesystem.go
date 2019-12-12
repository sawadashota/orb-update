package filesystem

import (
	"io"
	"os"

	"gopkg.in/src-d/go-billy.v4"
)

type Filesystem interface {
}

type Os struct{}

func NewOs() *Os {
	return new(Os)
}

func (o *Os) Reader(path string) (io.ReadCloser, error) {
	return os.OpenFile(path, os.O_RDONLY, 0666)
}

func (o *Os) OverWriter(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

type Memory struct {
	fs billy.Filesystem
}

func NewMemory(fs billy.Filesystem) *Memory {
	return &Memory{
		fs: fs,
	}
}

func (m *Memory) Reader(path string) (io.ReadCloser, error) {
	return m.fs.OpenFile(path, os.O_RDONLY, 0666)
}

func (m *Memory) OverWriter(path string) (io.WriteCloser, error) {
	return m.fs.Create(path)
}
