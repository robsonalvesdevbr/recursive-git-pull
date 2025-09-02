package finder

import (
	"os"
	"path/filepath"

	"github.com/robsonalvesdevbr/recursive-git-pull/pkg/types"
)

// FindRepositories recursively finds all Git repositories in the given path
func FindRepositories(rootPath string, includePatterns, excludePatterns []string) ([]*types.Repository, error) {
	var repositories []*types.Repository

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this is a .git directory
		if info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repoName := filepath.Base(repoPath)

			// Apply include/exclude patterns
			if shouldSkipRepository(repoName, includePatterns, excludePatterns) {
				return filepath.SkipDir
			}

			repo := &types.Repository{
				Path: repoPath,
				Name: repoName,
			}
			repositories = append(repositories, repo)

			// Skip walking into .git directory
			return filepath.SkipDir
		}

		return nil
	})

	return repositories, err
}

// shouldSkipRepository checks if a repository should be skipped based on patterns
func shouldSkipRepository(repoName string, includePatterns, excludePatterns []string) bool {
	// If include patterns are specified, repository must match at least one
	if len(includePatterns) > 0 {
		matched := false
		for _, pattern := range includePatterns {
			if matched, _ := filepath.Match(pattern, repoName); matched {
				matched = true
				break
			}
		}
		if !matched {
			return true
		}
	}

	// Check exclude patterns
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, repoName); matched {
			return true
		}
	}

	return false
}

// IsGitRepository checks if the given path is a Git repository
func IsGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if info, err := os.Stat(gitDir); err == nil {
		return info.IsDir()
	}
	return false
}

// GetRepositoryStatus returns a simple status of the repository
func GetRepositoryStatus(repoPath string) string {
	if !IsGitRepository(repoPath) {
		return "not-a-git-repo"
	}
	return "unknown"
}