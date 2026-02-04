package migrations

import (
	"database/sql"
	"fmt"
)

// MigrateUserImagesTable creates the user_images table for storing
// whiteboard PNG exports in user-level storage accessible by Claude.
func MigrateUserImagesTable(db *sql.DB) error {
	// Check if table already exists
	var tableExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM sqlite_master
		WHERE type='table' AND name='user_images'
	`).Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check user_images table: %w", err)
	}

	if tableExists {
		return nil
	}

	_, err = db.Exec(`
		CREATE TABLE user_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			image_type TEXT NOT NULL DEFAULT 'whiteboard',
			filename TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			content_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_user_images_user ON user_images(user_id);
		CREATE INDEX idx_user_images_type ON user_images(user_id, image_type);
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_images table: %w", err)
	}

	return nil
}
