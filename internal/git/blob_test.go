package git

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

const (
	testRepoPath   = "../../repos/Teagan.Kris/testrepository.git" // Adjust as needed
	testBlobHash   = "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"   // Empty file blob hash
	testCommitHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"   // Initial commit hash
	testFilePath   = "README.md"                                  // Adjust as needed
)

func repoExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestGetBlobSize(t *testing.T) {
	if !repoExists(testRepoPath) {
		t.Skipf("Test repo not found at %s, skipping", testRepoPath)
	}
	blobHash := plumbing.NewHash(testBlobHash)
	res, err := GetBlobSize(testRepoPath, blobHash)
	if err != nil {
		t.Fatalf("GetBlobSize failed: %v", err)
	}
	if res.Size < 0 {
		t.Errorf("Invalid blob size: %d", res.Size)
	}
}

func TestGetFileBlobSize(t *testing.T) {
	if !repoExists(testRepoPath) {
		t.Skipf("Test repo not found at %s, skipping", testRepoPath)
	}
	res, err := GetFileBlobSize(testRepoPath, testFilePath, testCommitHash)
	if err != nil {
		t.Fatalf("GetFileBlobSize failed: %v", err)
	}
	if res.Size < 0 {
		t.Errorf("Invalid file blob size: %d", res.Size)
	}
}

func TestReadBlob(t *testing.T) {
	if !repoExists(testRepoPath) {
		t.Skipf("Test repo not found at %s, skipping", testRepoPath)
	}
	blobHash := plumbing.NewHash(testBlobHash)
	blob, err := ReadBlob(testRepoPath, blobHash)
	if err != nil {
		t.Fatalf("ReadBlob failed: %v", err)
	}
	if blob.Size != int64(len(blob.Content)) {
		t.Errorf("Blob size (%d) does not match content length (%d)", blob.Size, len(blob.Content))
	}
}
