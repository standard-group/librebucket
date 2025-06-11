package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Blob is a git blob object with content and size
type Blob struct {
	ID      plumbing.Hash
	Content []byte
	Size    int64
}

// BlobSizeResult contains information about a blob
type BlobSizeResult struct {
	Hash plumbing.Hash // The blob hash
	Size int64         // Size in bytes
	Path string        // Path to the blob in the repository
}

// GetBlobSize returns the size of a blob in bytes given its hash
func GetBlobSize(repoPath string, blobHash plumbing.Hash) (*BlobSizeResult, error) {
	// Open the repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the blob object
	blob, err := r.BlobObject(blobHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob: %w", err)
	}

	return &BlobSizeResult{
		Hash: blobHash,
		Size: blob.Size,
	}, nil
}

// GetFileBlobSize returns the size of a file at a specific commit
func GetFileBlobSize(repoPath, filePath, commitHash string) (*BlobSizeResult, error) {
	// Open the repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the commit
	hash := plumbing.NewHash(commitHash)
	commit, err := r.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	// Get the file in the commit
	file, err := commit.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &BlobSizeResult{
		Hash: file.Blob.Hash,
		Size: file.Blob.Size,
		Path: filePath,
	}, nil
}

// ReadBlob reads the content and size of a blob by its hash from the repository
func ReadBlob(repoPath string, blobHash plumbing.Hash) (*Blob, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	blobObj, err := r.BlobObject(blobHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob: %w", err)
	}

	rd, err := blobObj.Reader()
	if err != nil {
		return nil, fmt.Errorf("failed to get blob reader: %w", err)
	}
	defer rd.Close()

	content := make([]byte, blobObj.Size)
	_, err = rd.Read(content)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to read blob content: %w", err)
	}

	return &Blob{
		ID:      blobHash,
		Content: content,
		Size:    blobObj.Size,
	}, nil
}
