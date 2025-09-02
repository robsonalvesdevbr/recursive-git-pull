package types

import "time"

// Repository represents a Git repository
type Repository struct {
	Path   string
	Name   string
	Status string
}

// Config holds configuration for the tool
type Config struct {
	RootPath         string
	Command          string
	Parallel         bool
	MaxWorkers       int
	Timeout          time.Duration
	IgnoreDirty      bool
	IncludePatterns  []string
	ExcludePatterns  []string
	Verbose          bool
	AllBranches      bool
	NoColor          bool
}

// ExecutionResult represents the result of command execution
type ExecutionResult struct {
	Repository *Repository
	Command    string
	Success    bool
	Output     string
	Error      string
	Duration   time.Duration
}