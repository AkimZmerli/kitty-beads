// Package whiteboard provides the whiteboard PNG export vertical slice.
package whiteboard

// SaveArgs represents arguments for saving whiteboard snapshot.
type SaveArgs struct {
	ImageData string `json:"image_data"` // Base64 encoded PNG data
}

// SaveResponse represents the result of a save operation.
type SaveResponse struct {
	Success bool   `json:"success"`
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

// UserImageSaveArgs represents arguments for saving to user storage.
type UserImageSaveArgs struct {
	ImageData string `json:"image_data"` // Base64 encoded PNG data
}

// UserImageSaveResponse represents the result of saving to user storage.
type UserImageSaveResponse struct {
	Success  bool   `json:"success"`
	Path     string `json:"path"`      // Full path to saved file
	Latest   string `json:"latest"`    // Path to latest.png symlink
	Message  string `json:"message,omitempty"`
	FileSize int64  `json:"file_size"`
}

// UserImageGetResponse represents metadata about a user's latest image.
type UserImageGetResponse struct {
	Exists   bool   `json:"exists"`
	Path     string `json:"path,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
	SavedAt  string `json:"saved_at,omitempty"`
}
