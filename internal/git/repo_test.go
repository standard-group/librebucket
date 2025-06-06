package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndMetaRepo(t *testing.T) {
	dir := t.TempDir()
	repoPath := filepath.Join(dir, "testrepo.git")
	owner := "alice"
	public := true

	err := CreateRepo(repoPath, owner, public)
	if err != nil {
		t.Fatalf("CreateRepo failed: %v", err)
	}

	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		t.Fatalf("LoadRepoMeta failed: %v", err)
	}
	if meta.Owner != owner || meta.Public != public {
		t.Errorf("RepoMeta mismatch: got %+v, want owner=%s, public=%v", meta, owner, public)
	}

	ownerStatus, err := IsRepoOwner(repoPath, owner)
	if err != nil {
		t.Fatalf("IsRepoOwner failed: %v", err)
	}
	if !ownerStatus {
		t.Errorf("IsRepoOwner failed: should be true for owner")
	}

	publicStatus, err := IsRepoPublic(repoPath)
	if err != nil {
		t.Fatalf("IsRepoPublic failed: %v", err)
	}
	if !publicStatus {
		t.Errorf("IsRepoPublic failed: should be true for public repo")
	}
}

func TestCloneRepo(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "src.git")
	dstPath := filepath.Join(dir, "dst.git")
	owner := "bob"
	public := false

	if err := CreateRepo(srcPath, owner, public); err != nil {
		t.Fatalf("CreateRepo failed: %v", err)
	}

	// Add a dummy file to srcPath to test clone
	dummyFile := filepath.Join(srcPath, "README.md")
	if err := os.WriteFile(dummyFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	if err := CloneRepo(srcPath, dstPath); err != nil {
		t.Fatalf("CloneRepo failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dstPath, "README.md")); err != nil {
		t.Errorf("Cloned repo missing README.md: %v", err)
	}
}
