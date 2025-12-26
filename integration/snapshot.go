// Snapshot testing utilities for integration tests
package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	packageDirOnce sync.Once
	packageDir     string
)

// getPackageDir returns the directory where this package's source files live
func getPackageDir() string {
	packageDirOnce.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		packageDir = filepath.Dir(filename)
	})
	return packageDir
}

// SnapshotConfig controls snapshot behavior
type SnapshotConfig struct {
	UpdateSnapshots bool // Set to true to update all snapshots (via -update flag)
}

var snapshotConfig = SnapshotConfig{
	UpdateSnapshots: os.Getenv("UPDATE_SNAPSHOTS") == "1",
}

// Snapshot represents a saved test snapshot
type Snapshot struct {
	Name     string            `json:"name"`
	Files    map[string]string `json:"files"`    // relPath -> content
	Metadata map[string]any    `json:"metadata"` // Additional test metadata
}

// SnapshotMatcher handles snapshot comparison and updates
type SnapshotMatcher struct {
	t            *testing.T
	snapshotDir  string
	snapshotName string
}

// NewSnapshotMatcher creates a new snapshot matcher for a test
func NewSnapshotMatcher(t *testing.T, testName string) *SnapshotMatcher {
	t.Helper()

	// Get the package source directory (where the test files live)
	sourceDir := getPackageDir()

	// Create snapshots directory in source tree
	snapshotDir := filepath.Join(sourceDir, "__snapshots__")
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		t.Fatalf("failed to create snapshot directory: %v", err)
	}

	// Sanitize test name for filename
	safeName := strings.ReplaceAll(testName, "/", "_")
	safeName = strings.ReplaceAll(safeName, " ", "_")

	return &SnapshotMatcher{
		t:            t,
		snapshotDir:  snapshotDir,
		snapshotName: safeName,
	}
}

// MatchSnapshot compares generated files against saved snapshot
func (sm *SnapshotMatcher) MatchSnapshot(sitePath string, filesToCheck []string) {
	sm.t.Helper()

	snapshot := sm.captureSnapshot(sitePath, filesToCheck)
	snapshotPath := filepath.Join(sm.snapshotDir, sm.snapshotName+".json")

	if snapshotConfig.UpdateSnapshots {
		sm.updateSnapshot(snapshotPath, snapshot)
		sm.t.Logf("✓ Updated snapshot: %s", sm.snapshotName)
		return
	}

	sm.compareSnapshot(snapshotPath, snapshot)
}

// MatchSnapshotWithMetadata compares snapshot with additional metadata
func (sm *SnapshotMatcher) MatchSnapshotWithMetadata(sitePath string, filesToCheck []string, metadata map[string]any) {
	sm.t.Helper()

	snapshot := sm.captureSnapshot(sitePath, filesToCheck)
	snapshot.Metadata = metadata
	snapshotPath := filepath.Join(sm.snapshotDir, sm.snapshotName+".json")

	if snapshotConfig.UpdateSnapshots {
		sm.updateSnapshot(snapshotPath, snapshot)
		sm.t.Logf("✓ Updated snapshot: %s", sm.snapshotName)
		return
	}

	sm.compareSnapshot(snapshotPath, snapshot)
}

// captureSnapshot captures current state of specified files
func (sm *SnapshotMatcher) captureSnapshot(sitePath string, filesToCheck []string) *Snapshot {
	sm.t.Helper()

	snapshot := &Snapshot{
		Name:     sm.snapshotName,
		Files:    make(map[string]string),
		Metadata: make(map[string]any),
	}

	publicDir := filepath.Join(sitePath, "public")

	for _, relPath := range filesToCheck {
		fullPath := filepath.Join(publicDir, relPath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				// File doesn't exist - record as empty
				snapshot.Files[relPath] = ""
			} else {
				sm.t.Fatalf("failed to read file %s: %v", relPath, err)
			}
		} else {
			snapshot.Files[relPath] = string(content)
		}
	}

	return snapshot
}

// updateSnapshot saves the current snapshot to disk
func (sm *SnapshotMatcher) updateSnapshot(snapshotPath string, snapshot *Snapshot) {
	sm.t.Helper()

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		sm.t.Fatalf("failed to marshal snapshot: %v", err)
	}

	// Get absolute path for debugging
	absPath, _ := filepath.Abs(snapshotPath)
	sm.t.Logf("Writing snapshot to: %s", absPath)

	if err := os.WriteFile(snapshotPath, data, 0644); err != nil {
		sm.t.Fatalf("failed to write snapshot to %s: %v", absPath, err)
	}
}

// compareSnapshot compares current state against saved snapshot
func (sm *SnapshotMatcher) compareSnapshot(snapshotPath string, current *Snapshot) {
	sm.t.Helper()

	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		if os.IsNotExist(err) {
			sm.t.Errorf(
				"Snapshot does not exist: %s\n\n"+
					"Run tests with UPDATE_SNAPSHOTS=1 to create it:\n"+
					"  UPDATE_SNAPSHOTS=1 go test ./integration -run %s",
				sm.snapshotName, sm.t.Name())
			return
		}
		sm.t.Fatalf("failed to read snapshot: %v", err)
	}

	var saved Snapshot
	if err := json.Unmarshal(data, &saved); err != nil {
		sm.t.Fatalf("failed to unmarshal snapshot: %v", err)
	}

	// Compare files
	sm.compareFiles(&saved, current)
}

// compareFiles compares two file sets
func (sm *SnapshotMatcher) compareFiles(saved, current *Snapshot) {
	sm.t.Helper()

	// Check for missing files
	for relPath := range saved.Files {
		if _, exists := current.Files[relPath]; !exists {
			sm.t.Errorf("File missing in current output: %s", relPath)
		}
	}

	// Check for extra files
	for relPath := range current.Files {
		if _, exists := saved.Files[relPath]; !exists {
			sm.t.Errorf("Unexpected file in current output: %s", relPath)
		}
	}

	// Compare file contents
	for relPath, savedContent := range saved.Files {
		currentContent, exists := current.Files[relPath]
		if !exists {
			continue // Already reported above
		}

		if savedContent != currentContent {
			sm.reportDifference(relPath, savedContent, currentContent)
		}
	}
}

// reportDifference reports a snapshot mismatch with helpful diff
func (sm *SnapshotMatcher) reportDifference(relPath, expected, actual string) {
	sm.t.Helper()

	// Generate content hashes for quick comparison
	expectedHash := sm.hash(expected)
	actualHash := sm.hash(actual)

	sm.t.Errorf(`Snapshot mismatch for file: %s

Expected hash: %s
Actual hash:   %s

To update this snapshot, run:
  UPDATE_SNAPSHOTS=1 go test ./integration -run %s

Expected length: %d bytes
Actual length:   %d bytes

First 500 chars of expected:
%s

First 500 chars of actual:
%s
`,
		relPath,
		expectedHash[:16],
		actualHash[:16],
		sm.t.Name(),
		len(expected),
		len(actual),
		sm.truncate(expected, 500),
		sm.truncate(actual, 500),
	)
}

// hash generates a SHA256 hash of content
func (sm *SnapshotMatcher) hash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

// truncate truncates a string to maxLen characters
func (sm *SnapshotMatcher) truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + fmt.Sprintf("... (%d more bytes)", len(s)-maxLen)
}

// SnapshotFiles is a convenience function for common snapshot testing
func SnapshotFiles(t *testing.T, testName string, sitePath string, files []string) {
	t.Helper()
	matcher := NewSnapshotMatcher(t, testName)
	matcher.MatchSnapshot(sitePath, files)
}

// SnapshotFilesWithMetadata is a convenience function with metadata
func SnapshotFilesWithMetadata(
	t *testing.T,
	testName string,
	sitePath string,
	files []string,
	metadata map[string]any,
) {
	t.Helper()
	matcher := NewSnapshotMatcher(t, testName)
	matcher.MatchSnapshotWithMetadata(sitePath, files, metadata)
}
