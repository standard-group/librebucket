package git

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

// RepoMeta holds metadata for a repository
type RepoMeta struct {
	Public     bool               `json:"public"`
	Owner      string             `json:"owner"`
	StarsCount int                `json:"stars_count"`
	LastCommit string             `json:"last_commit"` // Can be commit hash or timestamp
	ForksCount int                `json:"forks_count"`
	Languages  map[string]float64 `json:"languages"` // Map of language name to percent
	CreatedAt  time.Time          `json:"created_at"`
}

// Metadata is stored in repos/{username}/{reponame}.meta.json
const (
	metadataFile    = ".meta.json"
	safeRepoBaseDir = "repos"
)

func resolveSafePath(baseDir, relativePath string) (string, error) {
	absPath, err := filepath.Abs(filepath.Join(baseDir, relativePath))
	if err != nil {
		return "", err
	}
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absPath, baseAbs) {
		return "", fmt.Errorf("unsafe path: %s", absPath)
	}
	return absPath, nil
}

// SaveRepoMeta saves metadata for a repository
func SaveRepoMeta(repoPath string, meta RepoMeta) error {
	safeRepoPath, err := resolveSafePath(safeRepoBaseDir, repoPath)
	if err != nil {
		return fmt.Errorf("invalid repo path: %w", err)
	}
	metaFilePath := filepath.Join(safeRepoPath, metadataFile)
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal repo metadata: %w", err)
	}
	return ioutil.WriteFile(metaFilePath, data, 0644)
}

// LoadRepoMeta loads metadata for a repository
func LoadRepoMeta(repoPath string) (RepoMeta, error) {
	safeRepoPath, err := resolveSafePath(safeRepoBaseDir, repoPath)
	if err != nil {
		return RepoMeta{}, fmt.Errorf("invalid repo path: %w", err)
	}
	metaFilePath := filepath.Join(safeRepoPath, metadataFile)
	data, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		return RepoMeta{}, fmt.Errorf("failed to read repo metadata: %w", err)
	}

	var meta RepoMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return RepoMeta{}, fmt.Errorf("failed to unmarshal repo metadata: %w", err)
	}

	return meta, nil
}

// IsRepoOwner checks if the given username is the owner of the repo
func IsRepoOwner(repoPath, username string) (bool, error) {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return false, err
	}
	return meta.Owner == username, nil
}

// IsRepoPublic checks if the repo is public
func IsRepoPublic(repoPath string) (bool, error) {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return false, err
	}
	return meta.Public, nil
}

// CloneRepo clones a git repository to the specified directory
func CloneRepo(url, directory string) error {
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

// CreateRepo initializes a new git repository in the specified directory and saves metadata
func CreateRepo(repoPath, owner string, public bool) error {
	safeRepoPath, err := resolveSafePath(safeRepoBaseDir, repoPath)
	if err != nil {
		return fmt.Errorf("invalid repo path: %w", err)
	}

	_, err = git.PlainInit(safeRepoPath, true)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	meta := RepoMeta{
		Owner:      owner,
		Public:     public,
		StarsCount: 0,
		LastCommit: "",
		ForksCount: 0,
		Languages:  make(map[string]float64),
		CreatedAt:  time.Now(),
	}
	return SaveRepoMeta(repoPath, meta)
}

// UpdateStars updates the stars count for a repo
func UpdateStars(repoPath string, stars int) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.StarsCount = stars
	return SaveRepoMeta(repoPath, meta)
}

// UpdateForks updates the forks count for a repo
func UpdateForks(repoPath string, forks int) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.ForksCount = forks
	return SaveRepoMeta(repoPath, meta)
}

// UpdateLastCommit sets the last commit hash or timestamp
func UpdateLastCommit(repoPath, commit string) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.LastCommit = commit
	return SaveRepoMeta(repoPath, meta)
}

// UpdateLanguages sets the languages map for a repo
func UpdateLanguages(repoPath string, languages map[string]float64) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.Languages = languages
	return SaveRepoMeta(repoPath, meta)
}
