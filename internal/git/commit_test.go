package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func createTestRepoWithCommit(t *testing.T, dir string) (repoPath, commitHash string) {
	repoPath = filepath.Join(dir, "repo.git")
	if err := CreateRepo(repoPath, "testuser", true); err != nil {
		t.Fatalf("CreateRepo failed: %v", err)
	}
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("PlainOpen failed: %v", err)
	}
	wt, err := r.Worktree()
	if err != nil {
		t.Fatalf("Worktree failed: %v", err)
	}
	file := filepath.Join(repoPath, "file.txt")
	if err := os.WriteFile(file, []byte("hello world"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	_, err = wt.Add("file.txt")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	commit, err := wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	return repoPath, commit.String()
}

func TestGetCommitHistory(t *testing.T) {
	dir := t.TempDir()
	repoPath, _ := createTestRepoWithCommit(t, dir)
	commits, err := GetCommitHistory(repoPath, 10)
	if err != nil {
		t.Fatalf("GetCommitHistory failed: %v", err)
	}
	if len(commits) == 0 {
		t.Errorf("Expected at least 1 commit, got 0")
	}
}

func TestGetCommitByHash(t *testing.T) {
	dir := t.TempDir()
	repoPath, hash := createTestRepoWithCommit(t, dir)
	commit, err := GetCommitByHash(repoPath, hash)
	if err != nil {
		t.Fatalf("GetCommitByHash failed: %v", err)
	}
	if commit.Hash != hash {
		t.Errorf("Commit hash mismatch: got %s, want %s", commit.Hash, hash)
	}
}

func TestGetCommitChanges(t *testing.T) {
	dir := t.TempDir()
	repoPath, hash := createTestRepoWithCommit(t, dir)
	changes, err := GetCommitChanges(repoPath, hash)
	if err != nil {
		t.Fatalf("GetCommitChanges failed: %v", err)
	}
	if len(changes) == 0 {
		t.Errorf("Expected at least 1 file change, got 0")
	}
}
