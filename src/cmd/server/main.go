// Package main provides the kitty-beads HTTP server.
// This server exposes a REST API for the Beads issue tracker
// and serves the Spec Kitty-style dashboard frontend.
package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/types"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

//go:embed frontend/dist/*
var frontendFiles embed.FS

var (
	port      = flag.Int("port", 8080, "HTTP server port")
	beadsDir  = flag.String("beads-dir", ".beads", "Path to .beads directory")
	openBrowser = flag.Bool("open", false, "Open browser on start")
)

// Server holds the HTTP server state
type Server struct {
	store   *sqlite.SQLiteStorage
	rootDir string
}

// Lane represents a kanban lane with issues
type Lane struct {
	Name   string       `json:"name"`
	Issues []IssueCard  `json:"issues"`
	Count  int          `json:"count"`
}

// IssueCard is a simplified issue for the kanban board
type IssueCard struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Priority    int       `json:"priority"`
	Status      string    `json:"status"`
	Assignee    string    `json:"assignee,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Labels      []string  `json:"labels,omitempty"`
	IsBlocked   bool      `json:"is_blocked"`
	Blockers    []string  `json:"blockers,omitempty"`
}

// FeatureSummary represents a feature/epic in the feature list
type FeatureSummary struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	KanbanStats map[string]int    `json:"kanban_stats"`
	Artifacts   map[string]bool   `json:"artifacts"`
	Meta        map[string]string `json:"meta,omitempty"`
	IsLegacy    bool              `json:"is_legacy"`
}

// FeaturesResponse is the response for /api/features
type FeaturesResponse struct {
	Features       []FeatureSummary  `json:"features"`
	ProjectPath    string            `json:"project_path"`
	WorktreesRoot  *string           `json:"worktrees_root"`
	ActiveWorktree *string           `json:"active_worktree"`
	ActiveMission  map[string]string `json:"active_mission"`
}

// KanbanResponse is the response for /api/kanban/{id}
type KanbanResponse struct {
	Lanes         map[string][]IssueCard `json:"lanes"`
	IsLegacy      bool                   `json:"is_legacy"`
	UpgradeNeeded bool                   `json:"upgrade_needed"`
}

func main() {
	flag.Parse()

	// Find .beads directory
	beadsPath, err := findBeadsDir(*beadsDir)
	if err != nil {
		log.Fatalf("Could not find .beads directory: %v", err)
	}

	// Initialize SQLite storage
	dbPath := filepath.Join(beadsPath, "beads.db")
	ctx := context.Background()
	store, err := sqlite.New(ctx, dbPath)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer store.Close()

	rootDir := filepath.Dir(beadsPath)
	server := &Server{
		store:   store,
		rootDir: rootDir,
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/features", server.handleFeatures)
	mux.HandleFunc("/api/kanban/", server.handleKanban)
	mux.HandleFunc("/api/issues", server.handleIssues)
	mux.HandleFunc("/api/issues/", server.handleIssue)
	mux.HandleFunc("/api/ready", server.handleReadyWork)
	mux.HandleFunc("/api/health", server.handleHealth)
	mux.HandleFunc("/api/constitution", server.handleConstitution)
	mux.HandleFunc("/api/diagnostics", server.handleDiagnostics)
	mux.HandleFunc("/api/artifact/", server.handleArtifact)

	// WebSocket terminal (supports /api/terminal/{sessionId} for multi-tab)
	mux.HandleFunc("/api/terminal/", server.handleTerminal)
	mux.HandleFunc("/api/terminal", server.handleTerminal) // Also handle without trailing slash

	// Legacy static files (keep for backwards compatibility)
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Serve React frontend from embedded dist folder
	frontendFS, _ := fs.Sub(frontendFiles, "frontend/dist")
	mux.Handle("/assets/", http.FileServer(http.FS(frontendFS)))

	// SPA fallback - serve index.html for all non-API routes
	mux.HandleFunc("/", server.handleSPA(frontendFS))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}()

	log.Printf("Kitty-Beads server starting on http://localhost:%d", *port)
	log.Printf("Project: %s", rootDir)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func findBeadsDir(path string) (string, error) {
	// If absolute path, use it directly
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("beads directory not found: %s", path)
	}

	// Walk up from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		beadsPath := filepath.Join(dir, ".beads")
		if _, err := os.Stat(beadsPath); err == nil {
			return beadsPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find .beads directory")
}

// handleSPA returns an http.HandlerFunc that serves the React SPA.
// For static files that exist, it serves them directly.
// For all other routes, it serves index.html (SPA routing).
func (s *Server) handleSPA(frontendFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}

		// Try to serve the file directly
		file, err := frontendFS.Open(path)
		if err == nil {
			defer file.Close()
			stat, err := file.Stat()
			if err == nil && !stat.IsDir() {
				// File exists, serve it
				content, _ := fs.ReadFile(frontendFS, path)
				contentType := getContentType(path)
				w.Header().Set("Content-Type", contentType)
				w.Write(content)
				return
			}
		}

		// File doesn't exist or is a directory - serve index.html for SPA routing
		indexContent, err := fs.ReadFile(frontendFS, "index.html")
		if err != nil {
			http.Error(w, "Frontend not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexContent)
	}
}

// getContentType returns the MIME type based on file extension
func getContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".woff"):
		return "font/woff"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	default:
		return "application/octet-stream"
	}
}

// handleDashboard serves the legacy dashboard (kept for backwards compatibility)
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := templateFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Dashboard not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":       "ok",
		"project_path": s.rootDir,
		"timestamp":    time.Now().Format(time.RFC3339),
	}
	writeJSON(w, response)
}

func (s *Server) handleFeatures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all issues and group by parent (epics)
	issues, err := s.store.SearchIssues(ctx, "", types.IssueFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find epics (top-level issues or issues with type "epic")
	epics := make(map[string]*types.Issue)
	childrenByParent := make(map[string][]*types.Issue)

	for _, issue := range issues {
		if issue.IssueType == types.TypeEpic || !strings.Contains(issue.ID, ".") {
			// This is an epic or top-level issue
			epics[issue.ID] = issue
		} else {
			// This is a child - extract parent ID
			parts := strings.Split(issue.ID, ".")
			if len(parts) > 1 {
				parentID := strings.Join(parts[:len(parts)-1], ".")
				childrenByParent[parentID] = append(childrenByParent[parentID], issue)
			}
		}
	}

	// Build feature summaries
	var features []FeatureSummary
	for id, epic := range epics {
		children := childrenByParent[id]
		stats := computeKanbanStats(children)

		features = append(features, FeatureSummary{
			ID:          id,
			Name:        epic.Title,
			Path:        id,
			KanbanStats: stats,
			Artifacts: map[string]bool{
				"spec":       epic.Description != "",
				"plan":       epic.Design != "",
				"tasks":      len(children) > 0,
				"research":   false,
				"contracts":  false,
				"data_model": false,
				"checklists": epic.AcceptanceCriteria != "",
			},
			Meta: map[string]string{
				"mission": "software-dev",
			},
			IsLegacy: false,
		})
	}

	// Sort by ID descending (most recent first)
	sort.Slice(features, func(i, j int) bool {
		return features[i].ID > features[j].ID
	})

	response := FeaturesResponse{
		Features:    features,
		ProjectPath: s.rootDir,
		ActiveMission: map[string]string{
			"name":        "Kitty-Beads",
			"domain":      "software-dev",
			"version":     "1.0.0",
			"slug":        "kitty-beads",
			"description": "Beads issue tracker with Spec Kitty UI",
		},
	}

	writeJSON(w, response)
}

func (s *Server) handleKanban(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract feature ID from path: /api/kanban/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Feature ID required", http.StatusBadRequest)
		return
	}
	featureID := parts[3]

	// Get all issues under this feature
	issues, err := s.store.SearchIssues(ctx, "", types.IssueFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter to issues under this feature
	var featureIssues []*types.Issue
	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, featureID+".") || issue.ID == featureID {
			featureIssues = append(featureIssues, issue)
		}
	}

	// Group by status into lanes
	lanes := map[string][]IssueCard{
		"planned":    {},
		"doing":      {},
		"for_review": {},
		"done":       {},
	}

	for _, issue := range featureIssues {
		// Skip the epic itself
		if issue.ID == featureID {
			continue
		}

		// Check if blocked
		isBlocked, blockers, _ := s.store.IsBlocked(ctx, issue.ID)

		card := IssueCard{
			ID:          issue.ID,
			Title:       issue.Title,
			Priority:    issue.Priority,
			Status:      string(issue.Status),
			Assignee:    issue.Assignee,
			Description: truncate(issue.Description, 200),
			CreatedAt:   issue.CreatedAt,
			UpdatedAt:   issue.UpdatedAt,
			Labels:      issue.Labels,
			IsBlocked:   isBlocked,
			Blockers:    blockers,
		}

		lane := statusToLane(issue.Status)
		lanes[lane] = append(lanes[lane], card)
	}

	// Sort each lane by priority
	for _, cards := range lanes {
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Priority < cards[j].Priority
		})
	}

	response := KanbanResponse{
		Lanes:         lanes,
		IsLegacy:      false,
		UpgradeNeeded: false,
	}

	writeJSON(w, response)
}

func (s *Server) handleIssues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		issues, err := s.store.SearchIssues(ctx, "", types.IssueFilter{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, issues)

	case http.MethodPost:
		var issue types.Issue
		if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.store.CreateIssue(ctx, &issue, "kitty-beads"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, issue)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleIssue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract issue ID from path: /api/issues/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Issue ID required", http.StatusBadRequest)
		return
	}
	issueID := strings.Join(parts[3:], "/")

	switch r.Method {
	case http.MethodGet:
		issue, err := s.store.GetIssue(ctx, issueID)
		if err != nil {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}
		writeJSON(w, issue)

	case http.MethodPut:
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.store.UpdateIssue(ctx, issueID, updates, "kitty-beads"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		issue, _ := s.store.GetIssue(ctx, issueID)
		writeJSON(w, issue)

	case http.MethodDelete:
		if err := s.store.DeleteIssue(ctx, issueID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleReadyWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	issues, err := s.store.GetReadyWork(ctx, types.WorkFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, issues)
}

func (s *Server) handleConstitution(w http.ResponseWriter, r *http.Request) {
	// Look for CLAUDE.md or README.md
	paths := []string{
		filepath.Join(s.rootDir, "CLAUDE.md"),
		filepath.Join(s.rootDir, "README.md"),
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write(data)
			return
		}
	}

	http.Error(w, "Constitution not found", http.StatusNotFound)
}

func (s *Server) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := s.store.GetStatistics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": map[string]interface{}{
			"total_issues": stats.TotalIssues,
			"open_issues":  stats.OpenIssues,
			"project_path": s.rootDir,
		},
		"issues":          []interface{}{},
		"recommendations": []interface{}{},
	}

	writeJSON(w, response)
}

func (s *Server) handleArtifact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract from path: /api/artifact/{featureID}/{artifactName}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Feature ID and artifact name required", http.StatusBadRequest)
		return
	}
	featureID := parts[3]
	artifactName := parts[4]

	// Get the epic/feature issue
	issue, err := s.store.GetIssue(ctx, featureID)
	if err != nil {
		http.Error(w, "Feature not found", http.StatusNotFound)
		return
	}

	var content string
	exists := false

	switch artifactName {
	case "spec":
		content = issue.Description
		exists = content != ""
	case "plan":
		content = issue.Design
		exists = content != ""
	case "tasks":
		// Get child issues as task list
		issues, _ := s.store.SearchIssues(ctx, "", types.IssueFilter{})
		var taskLines []string
		for _, iss := range issues {
			if strings.HasPrefix(iss.ID, featureID+".") {
				status := "[ ]"
				if iss.Status == types.StatusClosed {
					status = "[x]"
				}
				taskLines = append(taskLines, fmt.Sprintf("- %s %s: %s", status, iss.ID, iss.Title))
			}
		}
		content = strings.Join(taskLines, "\n")
		exists = len(taskLines) > 0
	case "checklists":
		content = issue.AcceptanceCriteria
		exists = content != ""
	default:
		content = ""
		exists = false
	}

	response := map[string]interface{}{
		"exists":  exists,
		"content": content,
	}

	writeJSON(w, response)
}

// Helper functions

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(data)
}

func statusToLane(status types.Status) string {
	switch status {
	case types.StatusOpen:
		return "planned"
	case types.StatusInProgress:
		return "doing"
	case types.StatusBlocked:
		return "doing" // Show blocked in doing lane
	case types.StatusClosed:
		return "done"
	default:
		return "planned"
	}
}

func computeKanbanStats(issues []*types.Issue) map[string]int {
	stats := map[string]int{
		"planned":    0,
		"doing":      0,
		"for_review": 0,
		"done":       0,
	}

	for _, issue := range issues {
		lane := statusToLane(issue.Status)
		stats[lane]++
	}

	return stats
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
