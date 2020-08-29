package ioutil2

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testWrite(t, "/tmp/atomic-file-test1.txt")
	testWrite(t, "atomic-file-test2.txt")
	testWrite(t, filepath.Join(tmpDir, "atomic-file-test3.txt"))
}

func testWrite(t *testing.T, fname string) {
	err := WriteFileAtomic(fname, []byte("⚛write\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(fname); err != nil {
		t.Fatal(err)
	}
}
