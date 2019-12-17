package filesystem

import (
	"io"
	"os"

	"gopkg.in/src-d/go-billy.v4"
)

// Filesystem to operate git repository
type Filesystem interface {
	Reader(path string) (io.ReadCloser, error)
	OverWriter(path string) (io.WriteCloser, error)
}

// Os .
type Os struct{}

// NewOs .
func NewOs() *Os {
	return new(Os)
}

// Reader .
func (o *Os) Reader(path string) (io.ReadCloser, error) {
	return os.OpenFile(path, os.O_RDONLY, 0666)
}

// OverWriter .
func (o *Os) OverWriter(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Memory .
type Memory struct {
	fs billy.Filesystem
}

// NewMemory .
func NewMemory(fs billy.Filesystem) *Memory {
	return &Memory{
		fs: fs,
	}
}

// Reader .
func (m *Memory) Reader(path string) (io.ReadCloser, error) {
	return m.fs.OpenFile(path, os.O_RDONLY, 0666)
}

// OverWriter .
func (m *Memory) OverWriter(path string) (io.WriteCloser, error) {
	return m.fs.Create(path)
}
