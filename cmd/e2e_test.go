package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")
	xmlPath := "testdata/india.xml"

	// Index the test data
	cmd := exec.Command("../wikifind", "index", xmlPath, indexPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Index command failed: %v\nOutput: %s", err, output)
	}

	// Check if index files exist
	if _, err := os.Stat(filepath.Join(indexPath, "indexa.idx")); os.IsNotExist(err) {
		t.Error("Index file not created")
	}

	// Search for a term
	cmd = exec.Command("../wikifind", "search", indexPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer func() { _ = stdin.Close() }()
		_, _ = stdin.Write([]byte("apple\n"))
	}()
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Search command failed: %v\nOutput: %s", err, output)
	}

	// Check output contains results
	if len(output) == 0 {
		t.Error("No search output")
	}
}
