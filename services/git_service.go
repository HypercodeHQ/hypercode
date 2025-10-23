package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type TreeEntry struct {
	Type string // "tree" (folder) or "blob" (file)
	Name string
	Path string
	Mode string
}

type GitService interface {
	ListBranches(repoPath string) ([]string, error)
	GetDefaultBranch(repoPath string) (string, error)
	ListTree(repoPath, ref, path string) ([]TreeEntry, error)
	GetFileContent(repoPath, ref, path string) ([]byte, error)
	IsFile(repoPath, ref, path string) (bool, error)
}

type gitService struct {
	reposBasePath string
}

func NewGitService(reposBasePath string) GitService {
	return &gitService{
		reposBasePath: reposBasePath,
	}
}

// ListBranches returns all branch names in the repository
func (s *gitService) ListBranches(repoPath string) ([]string, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	cmd.Dir = absPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list branches: %w (output: %s)", err, out.String())
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	branches := strings.Split(output, "\n")
	return branches, nil
}

// GetDefaultBranch returns the default branch of the repository (HEAD)
func (s *gitService) GetDefaultBranch(repoPath string) (string, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	cmd.Dir = absPath

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// If symbolic-ref fails, try to get the first branch
		branches, err := s.ListBranches(repoPath)
		if err != nil {
			return "", err
		}
		if len(branches) > 0 {
			return branches[0], nil
		}
		return "", fmt.Errorf("no branches found")
	}

	return strings.TrimSpace(out.String()), nil
}

// ListTree returns the contents of a directory at the given ref and path
func (s *gitService) ListTree(repoPath, ref, path string) ([]TreeEntry, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the tree path
	treePath := ref + ":"
	if path != "" {
		treePath = ref + ":" + path
	}

	cmd := exec.Command("git", "ls-tree", treePath)
	cmd.Dir = absPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list tree: %w (output: %s)", err, out.String())
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return []TreeEntry{}, nil
	}

	lines := strings.Split(output, "\n")
	entries := make([]TreeEntry, 0, len(lines))

	for _, line := range lines {
		// Format: <mode> <type> <hash>\t<name>
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		mode := parts[0]
		entryType := parts[1]
		// hash := parts[2]
		name := strings.Join(parts[3:], " ")

		// Handle tab-separated names
		if idx := strings.Index(line, "\t"); idx > 0 {
			name = line[idx+1:]
		}

		entryPath := name
		if path != "" {
			entryPath = filepath.Join(path, name)
		}

		entries = append(entries, TreeEntry{
			Type: entryType,
			Name: name,
			Path: entryPath,
			Mode: mode,
		})
	}

	// Sort: folders first (tree), then files (blob)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type == "tree"
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries, nil
}

// GetFileContent returns the contents of a file at the given ref and path
func (s *gitService) GetFileContent(repoPath, ref, path string) ([]byte, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	blobPath := ref + ":" + path
	cmd := exec.Command("git", "show", blobPath)
	cmd.Dir = absPath

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get file content: %w (stderr: %s)", err, stderr.String())
	}

	return out.Bytes(), nil
}

// IsFile checks if the given path is a file (blob) or directory (tree)
func (s *gitService) IsFile(repoPath, ref, path string) (bool, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return false, err
	}

	// Use git cat-file to check the type
	objectPath := ref + ":" + path
	cmd := exec.Command("git", "cat-file", "-t", objectPath)
	cmd.Dir = absPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to check object type: %w", err)
	}

	objectType := strings.TrimSpace(out.String())
	return objectType == "blob", nil
}
