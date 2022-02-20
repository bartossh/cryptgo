package filesbuf

import (
	"fmt"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	cases := []struct {
		success bool
		path    string
	}{
		{true, "./"},
		{false, "./notexisting"},
	}

	for _, c := range cases {
		var file *os.File
		var err error

		a := Agent{}
		if c.success {
			file, err = os.CreateTemp(c.path, "temp*.*")
			if err != nil {
				t.Errorf("cannot create temporary file in path %s, %s", c.path, err)
			}
			defer file.Close()
			defer os.Remove(fmt.Sprintf("%s/%s", c.path, file.Name()))
			if _, err = a.GetReadCloser(fmt.Sprintf("%s/%s", c.path, file.Name())); err != nil {
				t.Errorf("expected to create io.ReadCloser, got err %s", err)
			}
			continue
		}
		if _, err = a.GetReadCloser(fmt.Sprintf("%s/%s", c.path, "not_existing.txt")); err == nil {
			t.Error("expected error, got nil")
		}

	}
}

func TestWrite(t *testing.T) {
	cases := []struct {
		success bool
		path    string
		name    string
	}{
		{false, "./", "temp*.*"},
		{true, "./", "new-file-123456789.txt"},
	}

	for _, c := range cases {
		var file *os.File
		var err error

		a := Agent{}
		if !c.success {
			file, err = os.CreateTemp(c.path, c.name)
			if err != nil {
				t.Errorf("cannot create temporary file in path %s, %s", c.path, err)
			}
			defer file.Close()
			defer os.Remove(fmt.Sprintf("%s/%s", c.path, file.Name()))
			if _, err = a.GetWriteCloser(fmt.Sprintf("%s/%s", c.path, file.Name())); err != nil {
				t.Error("expected to get an error")
			}
			continue
		}
		f, err := a.GetWriteCloser(fmt.Sprintf("%s/%s", c.path, c.name))
		if err != nil {
			t.Errorf("expected to create file %s, got %s", c.name, err)
		}
		defer f.Close()
		defer os.Remove(fmt.Sprintf("%s/%s", c.path, c.name))

	}
}
