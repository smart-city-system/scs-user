package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileInfo contains information about an uploaded file
type FileInfo struct {
	OriginalName string
	SavedName    string
	Size         int64
	ContentType  string
	Extension    string
	Path         string
}

// AllowedImageTypes defines the allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/bmp":  true,
}

// AllowedVideoTypes defines the allowed video MIME types
var AllowedVideoTypes = map[string]bool{
	"video/mp4":        true,
	"video/avi":        true,
	"video/mov":        true,
	"video/wmv":        true,
	"video/flv":        true,
	"video/webm":       true,
	"video/mkv":        true,
	"video/quicktime":  true,
}

// SaveUploadedFile saves an uploaded file to the specified directory
func SaveUploadedFile(file *multipart.FileHeader, uploadDir string) (*FileInfo, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	extension := filepath.Ext(file.Filename)
	uniqueName := generateUniqueFilename(extension)
	filePath := filepath.Join(uploadDir, uniqueName)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	size, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentTypeFromExtension(extension)
	}

	return &FileInfo{
		OriginalName: file.Filename,
		SavedName:    uniqueName,
		Size:         size,
		ContentType:  contentType,
		Extension:    extension,
		Path:         filePath,
	}, nil
}

// ValidateImageFile validates if the uploaded file is a valid image
func ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", file.Size, maxSize)
	}

	// Check content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentTypeFromExtension(filepath.Ext(file.Filename))
	}

	if !AllowedImageTypes[contentType] {
		return fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, getAllowedTypes(AllowedImageTypes))
	}

	return nil
}

// ValidateVideoFile validates if the uploaded file is a valid video
func ValidateVideoFile(file *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", file.Size, maxSize)
	}

	// Check content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentTypeFromExtension(filepath.Ext(file.Filename))
	}

	if !AllowedVideoTypes[contentType] {
		return fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, getAllowedTypes(AllowedVideoTypes))
	}

	return nil
}

// generateUniqueFilename generates a unique filename with timestamp and UUID
func generateUniqueFilename(extension string) string {
	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s%s", timestamp, uniqueID, extension)
}

// getContentTypeFromExtension returns content type based on file extension
func getContentTypeFromExtension(extension string) string {
	extension = strings.ToLower(extension)
	switch extension {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/avi"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/wmv"
	case ".flv":
		return "video/flv"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/mkv"
	default:
		return "application/octet-stream"
	}
}

// getAllowedTypes returns a slice of allowed content types
func getAllowedTypes(allowedTypes map[string]bool) []string {
	types := make([]string, 0, len(allowedTypes))
	for contentType := range allowedTypes {
		types = append(types, contentType)
	}
	return types
}

// DeleteFile safely deletes a file
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}
	
	return os.Remove(filePath)
}

// GetFileSize returns the size of a file
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
