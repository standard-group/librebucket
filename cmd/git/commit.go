// internal/git/commit.go
package git

import (
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit represents a git commit with metadata
type Commit struct {
	Hash           string
	Author         string
	AuthorEmail    string
	Message        string
	Committer      string
	CommitterEmail string
	AuthoredAt     time.Time
	CommittedAt    time.Time
	Parents        []string
	Files          []CommitFile
}

// CommitFile represents a file changed in a commit
type CommitFile struct {
	Path        string
	ChangeType  string // e.g., "Added", "Modified", "Deleted"
	ContentType string
}

// GetCommitHistory returns the commit history for a repository
func GetCommitHistory(repoPath string) ([]*Commit, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	headRef, err := r.Head()
	if err != nil {
		// If there's no HEAD, the repo is likely empty
		return []*Commit{}, nil
	}

	startCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	cIter, err := r.Log(&git.LogOptions{From: startCommit.Hash})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	var commits []*Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, &Commit{
			Hash:           c.Hash.String(),
			Author:         c.Author.Name,
			AuthorEmail:    c.Author.Email,
			Message:        c.Message,
			Committer:      c.Committer.Name,
			CommitterEmail: c.Committer.Email,
			AuthoredAt:     c.Author.When,
			CommittedAt:    c.Committer.When,
			Parents:        getCommitParents(c),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commit history: %w", err)
	}

	return commits, nil
}

// GetCommitByHash retrieves a specific commit by its hash
func GetCommitByHash(repoPath, hash string) (*Commit, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	h := plumbing.NewHash(hash)

	c, err := r.CommitObject(h)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	return &Commit{
		Hash:           c.Hash.String(),
		Author:         c.Author.Name,
		AuthorEmail:    c.Author.Email,
		Message:        c.Message,
		Committer:      c.Committer.Name,
		CommitterEmail: c.Committer.Email,
		AuthoredAt:     c.Author.When,
		CommittedAt:    c.Committer.When,
		Parents:        getCommitParents(c),
	}, nil
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
				Path:        f.Name,
				ChangeType:  "added",
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

			// TODO: Calculate additions and deletions
			// filePatch.Chunks() contains hunk information
			// Need to iterate through chunks and count lines based on type

			changes = append(changes, CommitFile{
				Path:        path,
				ChangeType:  changeType,
				ContentType: contentType,
				// Additions and Deletions would go here
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

// getCommitParents extracts parent hashes from a commit object
func getCommitParents(c *object.Commit) []string {
	var parents []string
	for _, p := range c.ParentHashes {
		parents = append(parents, p.String())
	}
	return parents
}
