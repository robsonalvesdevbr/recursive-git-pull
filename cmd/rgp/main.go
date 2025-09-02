package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/robsonalvesdevbr/recursive-git-pull/internal/config"
	"github.com/robsonalvesdevbr/recursive-git-pull/internal/finder"
	"github.com/robsonalvesdevbr/recursive-git-pull/internal/git"
	"github.com/robsonalvesdevbr/recursive-git-pull/pkg/types"
)

func main() {
	cfg := config.ParseFlags()

	if cfg.Verbose {
		fmt.Printf("Starting recursive git command execution...\n")
		fmt.Printf("Root path: %s\n", cfg.RootPath)
		fmt.Printf("Command: git %s\n", cfg.Command)
		fmt.Printf("Parallel: %t\n", cfg.Parallel)
		if cfg.Parallel {
			fmt.Printf("Max workers: %d\n", cfg.MaxWorkers)
		}
		fmt.Printf("Timeout: %v\n", cfg.Timeout)
		fmt.Println()
	}

	// Find all Git repositories
	repositories, err := finder.FindRepositories(cfg.RootPath, cfg.IncludePatterns, cfg.ExcludePatterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding repositories: %v\n", err)
		os.Exit(1)
	}

	if len(repositories) == 0 {
		fmt.Println("No Git repositories found in the specified path.")
		os.Exit(0)
	}

	fmt.Printf("Found %d repositories:\n", len(repositories))
	for _, repo := range repositories {
		fmt.Printf("  - %s (%s)\n", repo.Name, repo.Path)
	}
	fmt.Println()

	// Execute command on all repositories
	executor := git.NewExecutor(cfg)
	start := time.Now()
	
	results := executor.ExecuteCommandOnRepositories(repositories, cfg.Command)
	
	totalDuration := time.Since(start)

	// Print summary
	printSummary(results, totalDuration, cfg.Verbose)

	// Exit with error code if any command failed
	for _, result := range results {
		if !result.Success {
			os.Exit(1)
		}
	}
}

func printSummary(results []*types.ExecutionResult, totalDuration time.Duration, verbose bool) {
	successful := 0
	failed := 0

	// Sort results by repository name for consistent output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repository.Name < results[j].Repository.Name
	})

	fmt.Println("Summary:")
	fmt.Printf("========\n")

	for _, result := range results {
		if result.Success {
			successful++
			fmt.Printf("✓ %s (%v)\n", result.Repository.Name, result.Duration)
		} else {
			failed++
			fmt.Printf("✗ %s (%v)\n", result.Repository.Name, result.Duration)
			if result.Error != "" {
				fmt.Printf("  Error: %s\n", result.Error)
			}
		}

		if verbose && result.Output != "" {
			fmt.Printf("  Output: %s\n", result.Output)
		}
	}

	fmt.Printf("\nTotal: %d repositories processed in %v\n", len(results), totalDuration)
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		fmt.Printf("\nSome repositories failed. Check the errors above.\n")
	}
}