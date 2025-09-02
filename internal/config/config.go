package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robsonalvesdevbr/recursive-git-pull/pkg/types"
)

// ParseFlags parses command line flags and returns configuration
func ParseFlags() *types.Config {
	config := &types.Config{}

	flag.StringVar(&config.RootPath, "path", ".", "Root path to search for Git repositories")
	flag.StringVar(&config.Command, "command", "pull", "Git command to execute")
	flag.BoolVar(&config.Parallel, "parallel", true, "Execute commands in parallel")
	flag.IntVar(&config.MaxWorkers, "workers", 4, "Maximum number of parallel workers")
	
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "30s", "Timeout for each command")
	
	flag.BoolVar(&config.IgnoreDirty, "ignore-dirty", false, "Ignore repositories with uncommitted changes")
	
	var includeStr, excludeStr string
	flag.StringVar(&includeStr, "include", "", "Comma-separated patterns to include repositories")
	flag.StringVar(&excludeStr, "exclude", "", "Comma-separated patterns to exclude repositories")
	
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&config.AllBranches, "all-branches", false, "Pull all branches (only works with pull command)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		showHelp()
		os.Exit(0)
	}

	// Validate root path
	if info, err := os.Stat(config.RootPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid path '%s': %v\n", config.RootPath, err)
		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a directory\n", config.RootPath)
		os.Exit(1)
	}

	// Validate command
	if config.Command == "" {
		fmt.Fprintf(os.Stderr, "Error: Git command cannot be empty\n")
		os.Exit(1)
	}

	// Validate max workers
	if config.MaxWorkers <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Number of workers must be positive\n")
		os.Exit(1)
	}

	// Parse timeout
	if timeout, err := time.ParseDuration(timeoutStr); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid timeout format: %v\n", err)
		os.Exit(1)
	} else {
		config.Timeout = timeout
	}

	// Parse include/exclude patterns
	if includeStr != "" {
		config.IncludePatterns = strings.Split(includeStr, ",")
		for i, pattern := range config.IncludePatterns {
			config.IncludePatterns[i] = strings.TrimSpace(pattern)
		}
	}

	if excludeStr != "" {
		config.ExcludePatterns = strings.Split(excludeStr, ",")
		for i, pattern := range config.ExcludePatterns {
			config.ExcludePatterns[i] = strings.TrimSpace(pattern)
		}
	}

	return config
}

func showHelp() {
	fmt.Println("Recursive Git Pull - Execute Git commands recursively on multiple repositories")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  rgp [options]")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  rgp -path ./workspace -command pull")
	fmt.Println("  rgp -path ./projects -command status -parallel=false")
	fmt.Println("  rgp -path ./repos -command pull -all-branches")
	fmt.Println("  rgp -include '*-service' -exclude 'test-*'")
}