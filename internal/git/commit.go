// internal/git/commit.go
package git

import (
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit represents a git commit with metadata
type Commit struct {
	Hash         string    `json:"hash"`
	ShortHash    string    `json:"short_hash"`
	Author       string    `json:"author"`
	Email        string    `json:"email"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	ParentHashes []string  `json:"parent_hashes,omitempty"`
}

// CommitFile represents a file changed in a commit
// ContentType is set using the file extension
type CommitFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Type        string `json:"type"`        // "added", "modified", "deleted"
	ContentType string `json:"contentType"` // MIME type of the file
}

// GetCommitHistory returns the commit history for a repository
func GetCommitHistory(repoPath string, limit int) ([]Commit, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commit object for the HEAD reference
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	var commits []Commit
	count := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		if limit > 0 && count >= limit {
			return io.EOF // Stop after reaching the limit
		}

		// Get parent hashes
		parentHashes := []string{}
		for _, parent := range c.ParentHashes {
			parentHashes = append(parentHashes, parent.String())
		}

		commit := Commit{
			Hash:         c.Hash.String(),
			ShortHash:    c.Hash.String()[:7],
			Author:       c.Author.Name,
			Email:        c.Author.Email,
			Message:      c.Message,
			CreatedAt:    c.Author.When,
			ParentHashes: parentHashes,
		}
		commits = append(commits, commit)
		count++
		return nil
	})

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error iterating commits: %w", err)
	}

	return commits, nil
}

// GetCommitByHash retrieves a specific commit by its hash
func GetCommitByHash(repoPath, hash string) (*Commit, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Parse the hash
	h := plumbing.NewHash(hash)

	// Get the commit object
	c, err := r.CommitObject(h)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit %s: %w", hash, err)
	}

	// Get parent hashes
	parentHashes := []string{}
	for _, parent := range c.ParentHashes {
		parentHashes = append(parentHashes, parent.String())
	}

	commit := &Commit{
		Hash:         c.Hash.String(),
		ShortHash:    c.Hash.String()[:7],
		Author:       c.Author.Name,
		Email:        c.Author.Email,
		Message:      c.Message,
		CreatedAt:    c.Author.When,
		ParentHashes: parentHashes,
	}

	return commit, nil
}

// GetCommitChanges returns the files changed in a specific commit
func GetCommitChanges(repoPath, hash string) ([]CommitFile, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Parse the hash
	h := plumbing.NewHash(hash)

	// Get the commit object
	c, err := r.CommitObject(h)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit %s: %w", hash, err)
	}

	// Get parent commit to compare with
	var parentCommit *object.Commit
	if len(c.ParentHashes) > 0 {
		parentCommit, err = c.Parent(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent commit: %w", err)
		}
	}

	var changes []CommitFile
	if parentCommit == nil {
		// For the first commit, get all files
		fileIter, err := c.Files()
		if err != nil {
			return nil, fmt.Errorf("failed to get files: %w", err)
		}

		err = fileIter.ForEach(func(f *object.File) error {
			contentType := mime.TypeByExtension(filepath.Ext(f.Name))
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			changes = append(changes, CommitFile{
				Name:        f.Name,
				Path:        f.Name,
				Additions:   0, // No comparison available
				Deletions:   0, // No comparison available
				Type:        "added",
				ContentType: contentType,
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error iterating files: %w", err)
		}
	} else {
		// Get changes between this commit and parent
		patch, err := parentCommit.Patch(c)
		if err != nil {
			return nil, fmt.Errorf("failed to create patch: %w", err)
		}

		for _, filePatch := range patch.FilePatches() {
			from, to := filePatch.Files()

			var path string
			var changeType string
			var contentType string

			if from == nil {
				// File was added
				path = to.Path()
				changeType = "added"
				contentType = mime.TypeByExtension(filepath.Ext(path))
			} else if to == nil {
				// File was deleted
				path = from.Path()
				changeType = "deleted"
				contentType = mime.TypeByExtension(filepath.Ext(path))
			} else {
				// File was modified
				path = from.Path()
				changeType = "modified"
				contentType = mime.TypeByExtension(filepath.Ext(path))
			}
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			additions := 0
			deletions := 0
			for _, chunk := range filePatch.Chunks() {
				if chunk.Type() == diff.Add {
					additions += strings.Count(chunk.Content(), "\n")
					if !strings.HasSuffix(chunk.Content(), "\n") {
						additions++
					}
				} else if chunk.Type() == diff.Delete {
					deletions += strings.Count(chunk.Content(), "\n")
					if !strings.HasSuffix(chunk.Content(), "\n") {
						deletions++
					}
				}
			}

			changes = append(changes, CommitFile{
				Name:        path,
				Path:        path,
				Additions:   additions,
				Deletions:   deletions,
				Type:        changeType,
				ContentType: contentType,
			})
		}
	}

	return changes, nil
}

// GetFileAtCommit returns the content of a file at a specific commit
func GetFileAtCommit(repoPath, filePath, commitHash string) ([]byte, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Parse the hash
	h := plumbing.NewHash(commitHash)

	// Get the commit object
	c, err := r.CommitObject(h)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit %s: %w", commitHash, err)
	}

	// Get the file from the commit tree
	f, err := c.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("file %s not found in commit %s: %w", filePath, commitHash, err)
	}

	// Get the file contents
	content, err := f.Contents()
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents: %w", err)
	}

	return []byte(content), nil
}
