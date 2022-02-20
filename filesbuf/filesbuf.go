package filesbuf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Agent instance allows to get Reader and WriteCloser
type Agent struct{}

// GetReadCloser returns reader if file of given path exists
func (a Agent) GetReadCloser(path string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("agent cannot open file for given path %s, %s", path, err)
	}
	return f, nil
}

// GetWriteCloser returns write close if folder in which file has to be created exists and there is no other file of given name
func (a Agent) GetWriteCloser(path string) (io.WriteCloser, error) {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("agent cannot create file for given path %s, %s", path, err)
	}
	return f, nil
}
