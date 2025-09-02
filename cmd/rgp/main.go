package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/robsonalvesdevbr/recursive-git-pull/internal/colors"
	"github.com/robsonalvesdevbr/recursive-git-pull/internal/config"
	"github.com/robsonalvesdevbr/recursive-git-pull/internal/finder"
	"github.com/robsonalvesdevbr/recursive-git-pull/internal/git"
	"github.com/robsonalvesdevbr/recursive-git-pull/pkg/types"
)

func main() {
	cfg := config.ParseFlags()
	
	// Set color preferences
	if cfg.NoColor {
		colors.SetForceNoColor(true)
	}

	if cfg.Verbose {
		fmt.Printf("%s\n", colors.Bold("Starting recursive git command execution..."))
		fmt.Printf("%s %s\n", colors.Info("Root path:"), colors.Dim(cfg.RootPath))
		fmt.Printf("%s %s\n", colors.Info("Command:"), colors.Bold("git "+cfg.Command))
		fmt.Printf("%s %s\n", colors.Info("Parallel:"), colors.Bold(fmt.Sprintf("%t", cfg.Parallel)))
		if cfg.Parallel {
			fmt.Printf("%s %s\n", colors.Info("Max workers:"), colors.Bold(fmt.Sprintf("%d", cfg.MaxWorkers)))
		}
		fmt.Printf("%s %s\n", colors.Info("Timeout:"), colors.Bold(fmt.Sprintf("%v", cfg.Timeout)))
		fmt.Println()
	}

	// Find all Git repositories
	repositories, err := finder.FindRepositories(cfg.RootPath, cfg.IncludePatterns, cfg.ExcludePatterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s\n", colors.ErrorIcon(), colors.Error(fmt.Sprintf("Error finding repositories: %v", err)))
		os.Exit(1)
	}

	if len(repositories) == 0 {
		fmt.Printf("%s %s\n", colors.WarningIcon(), colors.Warning("No Git repositories found in the specified path."))
		os.Exit(0)
	}

	fmt.Printf("%s %s\n", colors.SuccessIcon(), colors.Success(fmt.Sprintf("Found %d repositories:", len(repositories))))
	for _, repo := range repositories {
		fmt.Printf("  %s %s %s\n", colors.Info("•"), colors.Bold(repo.Name), colors.Dim(fmt.Sprintf("(%s)", repo.Path)))
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

	fmt.Printf("%s\n", colors.Bold("Summary:"))
	fmt.Printf("%s\n", colors.Dim("========"))

	for _, result := range results {
		if result.Success {
			successful++
			duration := colors.Dim(fmt.Sprintf("(%v)", result.Duration))
			fmt.Printf("%s %s %s\n", colors.SuccessIcon(), colors.Success(result.Repository.Name), duration)
		} else {
			failed++
			duration := colors.Dim(fmt.Sprintf("(%v)", result.Duration))
			fmt.Printf("%s %s %s\n", colors.ErrorIcon(), colors.Error(result.Repository.Name), duration)
			if result.Error != "" {
				if strings.Contains(result.Error, "skipped") {
					fmt.Printf("  %s %s\n", colors.WarningIcon(), colors.Warning(result.Error))
				} else {
					fmt.Printf("  %s %s\n", colors.ErrorIcon(), colors.Error(result.Error))
				}
			}
		}

		if verbose && result.Output != "" {
			fmt.Printf("  %s %s\n", colors.InfoIcon(), colors.Dim(result.Output))
		}
	}

	totalInfo := fmt.Sprintf("Total: %d repositories processed in %v", len(results), totalDuration)
	successInfo := fmt.Sprintf("Successful: %d", successful)
	failedInfo := fmt.Sprintf("Failed: %d", failed)

	fmt.Printf("\n%s\n", colors.Bold(totalInfo))
	if successful > 0 {
		fmt.Printf("%s\n", colors.Success(successInfo))
	}
	if failed > 0 {
		fmt.Printf("%s\n", colors.Error(failedInfo))
		fmt.Printf("\n%s %s\n", colors.WarningIcon(), colors.Warning("Some repositories failed. Check the errors above."))
	}
}