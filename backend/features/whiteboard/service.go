// Package whiteboard provides the whiteboard PNG export vertical slice.
package whiteboard

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// maxImagesPerUser is the number of images to keep per user before cleanup
const maxImagesPerUser = 5

// Service provides business logic for whiteboard operations.
type Service struct {
	rootDir string
	db      *sql.DB
}

// NewService creates a new whiteboard service.
func NewService(rootDir string) *Service {
	return &Service{rootDir: rootDir}
}

// NewServiceWithDB creates a new whiteboard service with database support.
func NewServiceWithDB(rootDir string, db *sql.DB) *Service {
	return &Service{rootDir: rootDir, db: db}
}

// SaveSnapshot saves a base64 PNG image to .whiteboard/snapshot.png.
func (s *Service) SaveSnapshot(args SaveArgs) (*SaveResponse, error) {
	if args.ImageData == "" {
		return nil, fmt.Errorf("image_data is required")
	}

	// Strip data URL prefix if present (e.g., "data:image/png;base64,")
	imageData := args.ImageData
	if idx := strings.Index(imageData, ","); idx != -1 {
		imageData = imageData[idx+1:]
	}

	// Decode base64
	pngData, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 image data: %w", err)
	}

	// Validate PNG magic bytes
	if len(pngData) < 8 || string(pngData[:8]) != "\x89PNG\r\n\x1a\n" {
		return nil, fmt.Errorf("invalid PNG format")
	}

	// Create .whiteboard directory
	whiteboardDir := filepath.Join(s.rootDir, ".whiteboard")
	if err := os.MkdirAll(whiteboardDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .whiteboard directory: %w", err)
	}

	// Write snapshot.png
	snapshotPath := filepath.Join(whiteboardDir, "snapshot.png")
	if err := os.WriteFile(snapshotPath, pngData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write snapshot: %w", err)
	}

	return &SaveResponse{
		Success: true,
		Path:    snapshotPath,
		Message: "Whiteboard saved to .whiteboard/snapshot.png",
	}, nil
}

// SaveUserImage saves a PNG to user-level storage at ~/.beads/whiteboard/{user}/{timestamp}.png
func (s *Service) SaveUserImage(args UserImageSaveArgs) (*UserImageSaveResponse, error) {
	if args.ImageData == "" {
		return nil, fmt.Errorf("image_data is required")
	}

	// Get current user
	username := os.Getenv("USER")
	if username == "" {
		username = "default"
	}

	// Strip data URL prefix if present
	imageData := args.ImageData
	if idx := strings.Index(imageData, ","); idx != -1 {
		imageData = imageData[idx+1:]
	}

	// Decode base64
	pngData, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 image data: %w", err)
	}

	// Validate PNG magic bytes
	if len(pngData) < 8 || string(pngData[:8]) != "\x89PNG\r\n\x1a\n" {
		return nil, fmt.Errorf("invalid PNG format")
	}

	// Compute content hash
	hash := sha256.Sum256(pngData)
	contentHash := hex.EncodeToString(hash[:])

	// Create user directory: ~/.beads/whiteboard/{username}/
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	userDir := filepath.Join(homeDir, ".beads", "whiteboard", username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s.png", timestamp)
	filePath := filepath.Join(userDir, filename)

	// Write the PNG file
	if err := os.WriteFile(filePath, pngData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write image file: %w", err)
	}

	// Update/create latest.png symlink
	latestPath := filepath.Join(userDir, "latest.png")
	_ = os.Remove(latestPath) // Remove existing symlink if present
	if err := os.Symlink(filePath, latestPath); err != nil {
		// Symlink failed, try copying instead (Windows compatibility)
		if copyErr := os.WriteFile(latestPath, pngData, 0644); copyErr != nil {
			return nil, fmt.Errorf("failed to create latest.png: %w", copyErr)
		}
	}

	fileSize := int64(len(pngData))

	// Record in database if available
	if s.db != nil {
		_, err = s.db.Exec(`
			INSERT INTO user_images (user_id, image_type, filename, file_path, file_size, content_hash, created_at)
			VALUES (?, 'whiteboard', ?, ?, ?, ?, datetime('now'))
		`, username, filename, filePath, fileSize, contentHash)
		if err != nil {
			// Log but don't fail - file was saved successfully
			fmt.Printf("Warning: failed to record image in database: %v\n", err)
		}

		// Cleanup old images
		s.cleanupOldImages(username, userDir)
	}

	return &UserImageSaveResponse{
		Success:  true,
		Path:     filePath,
		Latest:   latestPath,
		Message:  fmt.Sprintf("Whiteboard saved to %s", latestPath),
		FileSize: fileSize,
	}, nil
}

// GetLatestUserImage returns metadata about the user's most recent whiteboard image.
func (s *Service) GetLatestUserImage() (*UserImageGetResponse, error) {
	username := os.Getenv("USER")
	if username == "" {
		username = "default"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	latestPath := filepath.Join(homeDir, ".beads", "whiteboard", username, "latest.png")

	info, err := os.Stat(latestPath)
	if os.IsNotExist(err) {
		return &UserImageGetResponse{Exists: false}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat latest image: %w", err)
	}

	return &UserImageGetResponse{
		Exists:   true,
		Path:     latestPath,
		FileSize: info.Size(),
		SavedAt:  info.ModTime().Format(time.RFC3339),
	}, nil
}

// cleanupOldImages removes old images keeping only the last maxImagesPerUser.
func (s *Service) cleanupOldImages(username, userDir string) {
	if s.db == nil {
		return
	}

	// Get all images for user ordered by creation time
	rows, err := s.db.Query(`
		SELECT id, file_path FROM user_images
		WHERE user_id = ? AND image_type = 'whiteboard'
		ORDER BY created_at DESC
	`, username)
	if err != nil {
		return
	}
	defer rows.Close()

	type imageRecord struct {
		id   int64
		path string
	}

	var images []imageRecord
	for rows.Next() {
		var img imageRecord
		if err := rows.Scan(&img.id, &img.path); err == nil {
			images = append(images, img)
		}
	}

	// Keep only the last maxImagesPerUser
	if len(images) > maxImagesPerUser {
		toDelete := images[maxImagesPerUser:]
		for _, img := range toDelete {
			// Delete file
			_ = os.Remove(img.path)
			// Delete database record
			_, _ = s.db.Exec("DELETE FROM user_images WHERE id = ?", img.id)
		}
	}

	// Also clean up orphaned files not in database
	s.cleanupOrphanedFiles(userDir)
}

// cleanupOrphanedFiles removes PNG files in the directory that aren't tracked in DB.
func (s *Service) cleanupOrphanedFiles(userDir string) {
	entries, err := os.ReadDir(userDir)
	if err != nil {
		return
	}

	// Get all tracked file paths from database
	rows, err := s.db.Query(`
		SELECT file_path FROM user_images WHERE image_type = 'whiteboard'
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	tracked := make(map[string]bool)
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err == nil {
			tracked[path] = true
		}
	}

	// Collect orphaned files (not latest.png and not tracked)
	var orphans []string
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "latest.png" {
			continue
		}
		fullPath := filepath.Join(userDir, entry.Name())
		if !tracked[fullPath] && strings.HasSuffix(entry.Name(), ".png") {
			orphans = append(orphans, fullPath)
		}
	}

	// Sort by name (which includes timestamp) and keep newest orphans if needed
	sort.Strings(orphans)
	if len(orphans) > 0 {
		// Delete all orphaned files
		for _, path := range orphans {
			_ = os.Remove(path)
		}
	}
}
