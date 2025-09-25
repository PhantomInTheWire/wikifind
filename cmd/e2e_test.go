package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndToEnd(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index")
	xmlPath := "testdata/india.xml"

	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "../wikifind", ".")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", buildOutput)

	// Index the test data
	cmd := exec.Command("../wikifind", "index", xmlPath, indexPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Index command failed: %s", output)

	// Check if index files exist
	indexFile := filepath.Join(indexPath, "indexa.idx")
	assert.FileExists(t, indexFile, "Index file should be created")

	// Search for a term
	cmd = exec.Command("../wikifind", "search", indexPath)
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	go func() {
		defer func() { _ = stdin.Close() }()
		_, _ = stdin.Write([]byte("apple\n"))
	}()
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Search command failed: %s", output)

	// Check output contains results
	assert.NotEmpty(t, output, "Search should produce output")
	assert.True(t, strings.Contains(string(output), "results"), "Output should contain 'results'")
}
