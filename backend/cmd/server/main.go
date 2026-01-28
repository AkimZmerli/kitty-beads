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

	"github.com/steveyegge/beads/features/comments"
	"github.com/steveyegge/beads/features/compaction"
	"github.com/steveyegge/beads/features/dependencies"
	"github.com/steveyegge/beads/features/epics"
	"github.com/steveyegge/beads/features/export"
	"github.com/steveyegge/beads/features/gates"
	"github.com/steveyegge/beads/features/issues"
	"github.com/steveyegge/beads/features/kanban"
	"github.com/steveyegge/beads/features/labels"
	"github.com/steveyegge/beads/features/statistics"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/shared/middleware"
	sharedstorage "github.com/steveyegge/beads/shared/storage"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed frontend-dist/*
var frontendFiles embed.FS

var (
	port     = flag.Int("port", 8080, "HTTP server port")
	beadsDir = flag.String("beads-dir", ".beads", "Path to .beads directory")
)

// Server holds the HTTP server state
type Server struct {
	store              *sqlite.SQLiteStorage
	rootDir            string
	issueHandler       *issues.HTTPHandler
	kanbanHandler      *kanban.HTTPHandler
	labelHandler       *labels.HTTPHandler
	commentHandler     *comments.HTTPHandler
	dependencyHandler  *dependencies.HTTPHandler
	statisticsHandler  *statistics.HTTPHandler
	epicsHandler       *epics.HTTPHandler
	gatesHandler       *gates.HTTPHandler
	compactionHandler  *compaction.HTTPHandler
	exportHandler      *export.HTTPHandler
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

	// Create storage adapter for vertical slice architecture
	adapter := sharedstorage.NewAdapter(store)

	// Create services for all vertical slices
	issueService := issues.NewService(
		adapter.Issues(),
		store,
		adapter.Labels(),
		adapter.Dependencies(),
	)
	kanbanService := kanban.NewService(adapter.Kanban(), adapter.Issues())
	labelService := labels.NewService(adapter.Labels())
	commentService := comments.NewService(adapter.Comments(), adapter.Events())
	dependencyService := dependencies.NewService(adapter.Dependencies())
	statisticsService := statistics.NewService(adapter.Statistics())
	epicsService := epics.NewService(adapter.Issues(), adapter.Kanban(), rootDir)
	gatesService := gates.NewService(adapter.Issues())
	compactionService := compaction.NewService(store)
	exportService := export.NewService(adapter.Export())

	// Create HTTP handlers for all vertical slices
	issueHandler := issues.NewHTTPHandler(issueService)
	kanbanHandler := kanban.NewHTTPHandler(kanbanService)
	labelHandler := labels.NewHTTPHandler(labelService)
	commentHandler := comments.NewHTTPHandler(commentService)
	dependencyHandler := dependencies.NewHTTPHandler(dependencyService)
	statisticsHandler := statistics.NewHTTPHandler(statisticsService)
	epicsHandler := epics.NewHTTPHandler(epicsService)
	gatesHandler := gates.NewHTTPHandler(gatesService)
	compactionHandler := compaction.NewHTTPHandler(compactionService)
	exportHandler := export.NewHTTPHandler(exportService)

	server := &Server{
		store:              store,
		rootDir:            rootDir,
		issueHandler:       issueHandler,
		kanbanHandler:      kanbanHandler,
		labelHandler:       labelHandler,
		commentHandler:     commentHandler,
		dependencyHandler:  dependencyHandler,
		statisticsHandler:  statisticsHandler,
		epicsHandler:       epicsHandler,
		gatesHandler:       gatesHandler,
		compactionHandler:  compactionHandler,
		exportHandler:      exportHandler,
	}

	mux := http.NewServeMux()

	// API routes - Features overview
	mux.HandleFunc("/api/features", server.handleFeatures)
	mux.HandleFunc("/api/health", server.handleHealth)
	mux.HandleFunc("/api/diagnostics", server.handleDiagnostics)
	mux.HandleFunc("/api/artifact/", server.handleArtifact)

	// Issue routes - using vertical slice handlers
	mux.HandleFunc("/api/issues/count", server.issueHandler.HandleCount)
	mux.HandleFunc("/api/issues", server.handleIssuesRouter)
	mux.HandleFunc("/api/issues/", server.handleIssueRouter)

	// Kanban routes - using vertical slice handlers
	mux.HandleFunc("/api/kanban/", server.kanbanHandler.HandleKanban)
	mux.HandleFunc("/api/ready", server.kanbanHandler.HandleReady)
	mux.HandleFunc("/api/blocked", server.kanbanHandler.HandleBlocked)
	mux.HandleFunc("/api/stale", server.kanbanHandler.HandleStale)
	mux.HandleFunc("/api/epics/status", server.kanbanHandler.HandleEpicStatus)

	// Label routes - using vertical slice handlers
	mux.HandleFunc("/api/labels/search", server.labelHandler.HandleGetByLabel)
	mux.HandleFunc("/api/labels", server.handleLabelsRouter)
	mux.HandleFunc("/api/labels/", server.labelHandler.HandleGet)

	// Comment routes - using vertical slice handlers
	mux.HandleFunc("/api/comments", server.commentHandler.HandleAdd)
	mux.HandleFunc("/api/comments/", server.commentHandler.HandleList)
	mux.HandleFunc("/api/events/", server.commentHandler.HandleListEvents)

	// Dependency routes - using vertical slice handlers
	mux.HandleFunc("/api/dependencies/cycles", server.dependencyHandler.HandleCycles)
	mux.HandleFunc("/api/dependencies/tree/", server.dependencyHandler.HandleTree)
	mux.HandleFunc("/api/dependencies", server.handleDependenciesRouter)
	mux.HandleFunc("/api/dependencies/", server.dependencyHandler.HandleGet)
	mux.HandleFunc("/api/dependents/", server.dependencyHandler.HandleGetDependents)

	// Statistics routes - using vertical slice handlers
	mux.HandleFunc("/api/stats", server.statisticsHandler.HandleGet)
	mux.HandleFunc("/api/stats/molecule/", server.statisticsHandler.HandleGetMoleculeProgress)

	// Epic routes - using vertical slice handlers
	mux.HandleFunc("/api/epics", server.epicsHandler.HandleFeatures)
	mux.HandleFunc("/api/epics/", server.handleEpicsRouter)

	// Gate routes - using vertical slice handlers
	mux.HandleFunc("/api/gates", server.handleGatesRouter)
	mux.HandleFunc("/api/gates/", server.handleGateRouter)

	// Compaction routes - using vertical slice handlers
	mux.HandleFunc("/api/compact/stats", server.compactionHandler.HandleStats)
	mux.HandleFunc("/api/compact", server.compactionHandler.HandleCompact)

	// Export routes - using vertical slice handlers
	mux.HandleFunc("/api/export", server.exportHandler.HandleExport)
	mux.HandleFunc("/api/import", server.exportHandler.HandleImport)
	mux.HandleFunc("/api/sync/status", server.exportHandler.HandleSyncStatus)

	// WebSocket terminal (supports /api/terminal/{sessionId} for multi-tab)
	mux.HandleFunc("/api/terminal/", server.handleTerminal)
	mux.HandleFunc("/api/terminal", server.handleTerminal) // Also handle without trailing slash

	// Legacy static files (keep for backwards compatibility)
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Serve React frontend from embedded dist folder
	frontendFS, _ := fs.Sub(frontendFiles, "frontend-dist")
	mux.Handle("/assets/", http.FileServer(http.FS(frontendFS)))

	// SPA fallback - serve index.html for all non-API routes
	mux.HandleFunc("/", server.handleSPA(frontendFS))

	// Apply middleware chain
	handler := middleware.Chain(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
	)(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
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

// handleIssuesRouter routes /api/issues requests to the appropriate vertical slice handler
func (s *Server) handleIssuesRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.issueHandler.HandleList(w, r)
	case http.MethodPost:
		s.issueHandler.HandleCreate(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleIssueRouter routes /api/issues/{id} requests to the appropriate vertical slice handler
func (s *Server) handleIssueRouter(w http.ResponseWriter, r *http.Request) {
	// Check if this is a close operation: /api/issues/{id}/close
	if strings.HasSuffix(r.URL.Path, "/close") {
		s.issueHandler.HandleClose(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.issueHandler.HandleGet(w, r)
	case http.MethodPut, http.MethodPatch:
		s.issueHandler.HandleUpdate(w, r)
	case http.MethodDelete:
		s.issueHandler.HandleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLabelsRouter routes /api/labels requests to the appropriate vertical slice handler
func (s *Server) handleLabelsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.labelHandler.HandleAdd(w, r)
	case http.MethodDelete:
		s.labelHandler.HandleRemove(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDependenciesRouter routes /api/dependencies requests to the appropriate vertical slice handler
func (s *Server) handleDependenciesRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.dependencyHandler.HandleAdd(w, r)
	case http.MethodDelete:
		s.dependencyHandler.HandleRemove(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleEpicsRouter routes /api/epics/{id} requests to the appropriate vertical slice handler
func (s *Server) handleEpicsRouter(w http.ResponseWriter, r *http.Request) {
	// Check for status endpoint: /api/epics/status
	if strings.HasSuffix(r.URL.Path, "/status") || strings.Contains(r.URL.Path, "/status?") {
		s.epicsHandler.HandleStatus(w, r)
		return
	}

	// Otherwise, get specific epic summary: /api/epics/{id}
	switch r.Method {
	case http.MethodGet:
		s.epicsHandler.HandleGetSummary(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGatesRouter routes /api/gates requests to the appropriate vertical slice handler
func (s *Server) handleGatesRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.gatesHandler.HandleList(w, r)
	case http.MethodPost:
		s.gatesHandler.HandleCreate(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGateRouter routes /api/gates/{id} requests to the appropriate vertical slice handler
func (s *Server) handleGateRouter(w http.ResponseWriter, r *http.Request) {
	// Check for close operation: /api/gates/{id}/close
	if strings.HasSuffix(r.URL.Path, "/close") {
		s.gatesHandler.HandleClose(w, r)
		return
	}

	// Check for wait operation: /api/gates/{id}/wait
	if strings.HasSuffix(r.URL.Path, "/wait") {
		s.gatesHandler.HandleWait(w, r)
		return
	}

	// Otherwise, show gate: /api/gates/{id}
	switch r.Method {
	case http.MethodGet:
		s.gatesHandler.HandleShow(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
