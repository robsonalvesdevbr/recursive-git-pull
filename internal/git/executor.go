package git

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/robsonalvesdevbr/recursive-git-pull/pkg/types"
)

// Executor handles Git command execution
type Executor struct {
	config *types.Config
}

// NewExecutor creates a new Git executor
func NewExecutor(config *types.Config) *Executor {
	return &Executor{config: config}
}

// ExecuteCommand executes a Git command in a single repository
func (e *Executor) ExecuteCommand(repo *types.Repository, command string) *types.ExecutionResult {
	start := time.Now()
	result := &types.ExecutionResult{
		Repository: repo,
		Command:    command,
		Success:    false,
		Duration:   0,
	}

	// Check if we should ignore dirty repositories
	if e.config.IgnoreDirty && command == "pull" {
		if isDirty, err := e.isRepositoryDirty(repo.Path); err != nil {
			result.Error = fmt.Sprintf("Error checking repository status: %v", err)
			result.Duration = time.Since(start)
			return result
		} else if isDirty {
			result.Error = "Repository has uncommitted changes (skipped)"
			result.Duration = time.Since(start)
			return result
		}
	}

	// Handle special case for pull all branches
	if command == "pull" && e.config.AllBranches {
		return e.pullAllBranches(repo, start)
	}

	// Execute the command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", strings.Fields(command)...)
	cmd.Dir = repo.Path

	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Sprintf("Command timed out after %v", e.config.Timeout)
		} else {
			result.Error = err.Error()
		}
		result.Success = false
	} else {
		result.Success = true
	}

	return result
}

// ExecuteCommandOnRepositories executes a command on multiple repositories
func (e *Executor) ExecuteCommandOnRepositories(repositories []*types.Repository, command string) []*types.ExecutionResult {
	if !e.config.Parallel {
		return e.executeSequentially(repositories, command)
	}
	return e.executeInParallel(repositories, command)
}

// executeSequentially executes commands one by one
func (e *Executor) executeSequentially(repositories []*types.Repository, command string) []*types.ExecutionResult {
	results := make([]*types.ExecutionResult, 0, len(repositories))

	for _, repo := range repositories {
		if e.config.Verbose {
			fmt.Printf("Executing 'git %s' in %s...\n", command, repo.Path)
		}
		
		result := e.ExecuteCommand(repo, command)
		results = append(results, result)
		
		if e.config.Verbose {
			e.printResult(result)
		}
	}

	return results
}

// executeInParallel executes commands in parallel with worker pool
func (e *Executor) executeInParallel(repositories []*types.Repository, command string) []*types.ExecutionResult {
	jobsCh := make(chan *types.Repository, len(repositories))
	resultsCh := make(chan *types.ExecutionResult, len(repositories))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < e.config.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobsCh {
				if e.config.Verbose {
					fmt.Printf("Executing 'git %s' in %s...\n", command, repo.Path)
				}
				
				result := e.ExecuteCommand(repo, command)
				resultsCh <- result
				
				if e.config.Verbose {
					e.printResult(result)
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, repo := range repositories {
			jobsCh <- repo
		}
		close(jobsCh)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results
	results := make([]*types.ExecutionResult, 0, len(repositories))
	for result := range resultsCh {
		results = append(results, result)
	}

	return results
}

// isRepositoryDirty checks if repository has uncommitted changes
func (e *Executor) isRepositoryDirty(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

// pullAllBranches pulls all branches in the repository
func (e *Executor) pullAllBranches(repo *types.Repository, start time.Time) *types.ExecutionResult {
	result := &types.ExecutionResult{
		Repository: repo,
		Command:    "pull --all",
		Success:    false,
	}

	// Get all remote branches
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "branch", "-r")
	cmd.Dir = repo.Path

	output, err := cmd.Output()
	if err != nil {
		result.Error = fmt.Sprintf("Error getting remote branches: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	branches := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.Contains(line, "->") {
			// Remove "origin/" prefix
			branch := strings.TrimPrefix(line, "origin/")
			branches = append(branches, branch)
		}
	}

	// Pull each branch
	var outputs []string
	for _, branch := range branches {
		ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
		cmd := exec.CommandContext(ctx, "git", "pull", "origin", branch)
		cmd.Dir = repo.Path

		branchOutput, err := cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("Branch %s: %s", branch, string(branchOutput)))
		
		cancel()
		
		if err != nil {
			result.Error = fmt.Sprintf("Error pulling branch %s: %v", branch, err)
			result.Output = strings.Join(outputs, "\n")
			result.Duration = time.Since(start)
			return result
		}
	}

	result.Success = true
	result.Output = strings.Join(outputs, "\n")
	result.Duration = time.Since(start)
	return result
}

// printResult prints the execution result
func (e *Executor) printResult(result *types.ExecutionResult) {
	status := "✓"
	if !result.Success {
		status = "✗"
	}
	
	fmt.Printf("%s %s (%v)\n", status, result.Repository.Name, result.Duration)
	
	if result.Error != "" {
		fmt.Printf("  Error: %s\n", result.Error)
	}
	
	if result.Output != "" && e.config.Verbose {
		fmt.Printf("  Output: %s\n", strings.TrimSpace(result.Output))
	}
}