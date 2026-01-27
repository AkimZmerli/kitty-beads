package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsValidRemoteURL tests the remote URL validation function
func TestIsValidRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		// Valid URLs
		{"dolthub scheme", "dolthub://org/repo", true},
		{"gs scheme", "gs://bucket/path", true},
		{"s3 scheme", "s3://bucket/path", true},
		{"file scheme", "file:///path/to/repo", true},
		{"https scheme", "https://github.com/user/repo", true},
		{"http scheme", "http://github.com/user/repo", true},
		{"ssh scheme", "ssh://git@github.com/user/repo", true},
		{"git ssh format", "git@github.com:user/repo.git", true},
		{"git ssh with underscore", "git@gitlab.example_host.com:user/repo.git", true},

		// Invalid URLs
		{"empty string", "", false},
		{"no scheme", "github.com/user/repo", false},
		{"invalid scheme", "ftp://server/path", false},
		{"malformed git ssh", "git@:repo", false},
		{"just path", "/path/to/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRemoteURL(tt.url)
			if got != tt.expected {
				t.Errorf("isValidRemoteURL(%q) = %v, want %v", tt.url, got, tt.expected)
			}
		})
	}
}

// TestValidateSyncConfig tests the sync config validation function
func TestValidateSyncConfig(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("Failed to create .beads dir: %v", err)
	}

	t.Run("valid empty config", func(t *testing.T) {
		// Create minimal config.yaml
		configContent := `prefix: test
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		if len(issues) != 0 {
			t.Errorf("Expected no issues for valid empty config, got: %v", issues)
		}
	})

	t.Run("invalid sync.mode", func(t *testing.T) {
		configContent := `prefix: test
sync:
  mode: "invalid-mode"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		found := false
		for _, issue := range issues {
			if strings.Contains(issue, "sync.mode") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue about sync.mode, got: %v", issues)
		}
	})

	t.Run("invalid conflict.strategy", func(t *testing.T) {
		configContent := `prefix: test
conflict:
  strategy: "invalid-strategy"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		found := false
		for _, issue := range issues {
			if strings.Contains(issue, "conflict.strategy") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue about conflict.strategy, got: %v", issues)
		}
	})

	t.Run("invalid federation.sovereignty", func(t *testing.T) {
		configContent := `prefix: test
federation:
  sovereignty: "invalid-value"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		found := false
		for _, issue := range issues {
			if strings.Contains(issue, "federation.sovereignty") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue about federation.sovereignty, got: %v", issues)
		}
	})

	t.Run("external mode without remote", func(t *testing.T) {
		configContent := `prefix: test
sync:
  mode: "external"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		found := false
		for _, issue := range issues {
			if strings.Contains(issue, "federation.remote") && strings.Contains(issue, "required") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue about federation.remote being required, got: %v", issues)
		}
	})

	t.Run("invalid remote URL", func(t *testing.T) {
		configContent := `prefix: test
federation:
  remote: "invalid-url"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		found := false
		for _, issue := range issues {
			if strings.Contains(issue, "federation.remote") && strings.Contains(issue, "not a valid remote URL") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected issue about invalid remote URL, got: %v", issues)
		}
	})

	t.Run("valid sync config", func(t *testing.T) {
		configContent := `prefix: test
sync:
  mode: "git-branch"
conflict:
  strategy: "lww"
federation:
  sovereignty: "federated"
  remote: "https://github.com/user/beads-data.git"
`
		if err := os.WriteFile(filepath.Join(beadsDir, "config.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config.yaml: %v", err)
		}

		issues := validateSyncConfig(tmpDir)
		if len(issues) != 0 {
			t.Errorf("Expected no issues for valid config, got: %v", issues)
		}
	})
}

// TestFindBeadsRepoRoot tests the repo root finding function
func TestFindBeadsRepoRoot(t *testing.T) {
	// Create a temp directory structure
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	subDir := filepath.Join(tmpDir, "sub", "dir")

	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("Failed to create .beads dir: %v", err)
	}
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub dir: %v", err)
	}

	t.Run("from repo root", func(t *testing.T) {
		got := findBeadsRepoRoot(tmpDir)
		if got != tmpDir {
			t.Errorf("findBeadsRepoRoot(%q) = %q, want %q", tmpDir, got, tmpDir)
		}
	})

	t.Run("from subdirectory", func(t *testing.T) {
		got := findBeadsRepoRoot(subDir)
		if got != tmpDir {
			t.Errorf("findBeadsRepoRoot(%q) = %q, want %q", subDir, got, tmpDir)
		}
	})

	t.Run("not in repo", func(t *testing.T) {
		noRepoDir := t.TempDir()
		got := findBeadsRepoRoot(noRepoDir)
		if got != "" {
			t.Errorf("findBeadsRepoRoot(%q) = %q, want empty string", noRepoDir, got)
		}
	})
}
