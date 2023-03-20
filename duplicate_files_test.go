package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// create a temporary directory
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// create some files with the same content
	content := "hello world"
	files := []string{
		"file1.txt",
		"dir/file2.txt",
		"dir/subdir/file3.txt",
		"dir/subdir/file4.txt",
	}
	for _, file := range files {
		path := filepath.Join(dir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// run the script
	hashes := run(dir)

	// check the results
	if len(hashes) != 1 {
		t.Fatalf("expected 1 hash, got %d", len(hashes))
	}
	var hash string
	for h := range hashes {
		hash = h
		break
	}
	expectedHash := "5eb63bbbe01eeed093cb22bb8f5acdc3"
	if hash != expectedHash {
		t.Fatalf("expected hash %s, got %s", expectedHash, hash)
	}
	filesWithHash := hashes[hash]
	if len(filesWithHash) != len(files) {
		t.Fatalf("expected %d files with hash, got %d", len(files), len(filesWithHash))
	}
	for _, file := range files {
		found := false
		for _, f := range filesWithHash {
			if strings.HasSuffix(f, file) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected file %s to have hash %s", file, hash)
		}
	}
}

