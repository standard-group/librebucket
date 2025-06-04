package git

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
)

// RepoMeta holds metadata for a repository
// Public: true if public, false if private
// Owner: username of the owner
// StarsCount: number of stars
// LastCommit: last commit hash or timestamp
// ForksCount: number of forks
// Languages: map of language name to percent
// CreatedAt: repo creation time
//
// Metadata is stored in repos/{username}/{reponame}.meta.json
type RepoMeta struct {
	Owner      string             `json:"owner"`
	Public     bool               `json:"public"`
	StarsCount int                `json:"stars_count"`
	LastCommit string             `json:"last_commit"`
	ForksCount int                `json:"forks_count"`
	Languages  map[string]float64 `json:"languages"`
	CreatedAt  time.Time          `json:"created_at"`
}

// SaveRepoMeta saves metadata for a repository
func SaveRepoMeta(repoPath string, meta RepoMeta) error {
	metaPath := repoPath + ".meta.json"
	f, err := os.Create(metaPath)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(meta)
}

// LoadRepoMeta loads metadata for a repository
func LoadRepoMeta(repoPath string) (RepoMeta, error) {
	metaPath := repoPath + ".meta.json"
	f, err := os.Open(metaPath)
	if err != nil {
		return RepoMeta{}, fmt.Errorf("failed to open meta file: %w", err)
	}
	defer f.Close()
	var meta RepoMeta
	if err := json.NewDecoder(f).Decode(&meta); err != nil {
		return RepoMeta{}, fmt.Errorf("failed to decode meta file: %w", err)
	}
	return meta, nil
}

// IsRepoOwner checks if the given username is the owner of the repo
func IsRepoOwner(repoPath, username string) bool {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return false
	}
	return meta.Owner == username
}

// IsRepoPublic checks if the repo is public
func IsRepoPublic(repoPath string) bool {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return false
	}
	return meta.Public
}

// CloneRepo clones a git repository to the specified directory
func CloneRepo(url, directory string) error {
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: nil,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}
	return nil
}

// CreateRepo initializes a new git repository in the specified directory and saves metadata
// CreateRepo initializes a new bare git repository and saves metadata
func CreateRepo(directory, owner string, public bool) error {
	// Create as bare repository for server-side hosting
	_, err := git.PlainInit(directory, true) // true = bare repository
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	meta := RepoMeta{
		Owner:      owner,
		Public:     public,
		StarsCount: 0,
		LastCommit: "",
		ForksCount: 0,
		Languages:  map[string]float64{},
		CreatedAt:  time.Now(),
	}
	return SaveRepoMeta(directory, meta)
}

// UpdateStars updates the stars count for a repo
func UpdateStars(repoPath string, delta int) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.StarsCount += delta
	if meta.StarsCount < 0 {
		meta.StarsCount = 0
	}
	return SaveRepoMeta(repoPath, meta)
}

// UpdateForks updates the forks count for a repo
func UpdateForks(repoPath string, delta int) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.ForksCount += delta
	if meta.ForksCount < 0 {
		meta.ForksCount = 0
	}
	return SaveRepoMeta(repoPath, meta)
}

// UpdateLastCommit sets the last commit hash or timestamp
func UpdateLastCommit(repoPath, lastCommit string) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.LastCommit = lastCommit
	return SaveRepoMeta(repoPath, meta)
}

// UpdateLanguages sets the languages map for a repo
func UpdateLanguages(repoPath string, langs map[string]float64) error {
	meta, err := LoadRepoMeta(repoPath)
	if err != nil {
		return err
	}
	meta.Languages = langs
	return SaveRepoMeta(repoPath, meta)
}
