package ioutil2

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Write file to temp and atomically move when everything else succeeds.
func WriteFileAtomic(filename string, data []byte, perm os.FileMode) (err error) {
	wr, err := NewAtomicFileWriter(filename, perm)
	if err != nil {
		return err
	}

	_, err = wr.Write(data)
	// If no error already, propagate one.
	if closeErr := wr.Close(); err == nil {
		err = closeErr
	}

	return err
}

type atomicFileWriter struct {
	dstFilename string
	f           *os.File
}

func (afw *atomicFileWriter) Write(data []byte) (n int, err error) {
	return afw.f.Write(data)
}

func (afw *atomicFileWriter) Close() (err error) {
	err = afw.f.Sync()
	if closeErr := afw.f.Close(); err == nil {
		err = closeErr
	}
	// Any err should result in full cleanup.
	if err != nil {
		_ = os.Remove(afw.f.Name())
		return err
	}

	if err = os.Rename(afw.f.Name(), afw.dstFilename); err != nil {
		return err
	}

	fDir, err := os.Open(filepath.Dir(afw.dstFilename))
	if err != nil {
		return err
	}
	defer func() {
		// If no error already, propagate one.
		if closeErr := fDir.Close(); err == nil {
			err = closeErr
		}
	}()

	return fDir.Sync()
}

// Return a WriteCloser that accumulates data in a tempfile and atomically moves to the final destination when closed.
func NewAtomicFileWriter(filename string, perm os.FileMode) (wr io.WriteCloser, err error) {
	dir := filepath.Dir(filename)
	name := filepath.Base(filename)

	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(f.Name(), perm); err != nil {
		return nil, err
	}

	return &atomicFileWriter{dstFilename: filename, f: f}, nil
}
