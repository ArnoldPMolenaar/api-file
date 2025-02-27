package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

// Base64ToBytes func for convert base64 string to bytes.
func Base64ToBytes(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

// GetMimeTypeAndBase64 extracts the mimetype and data from a base64 string.
func GetMimeTypeAndBase64(value string) (string, string, error) {
	if idx := strings.Index(value, ";base64,"); idx != -1 {
		mimeType := value[:idx]
		data := value[idx+8:]

		return mimeType, data, nil
	}
	return "", "", errors.New("invalid base64 string")
}

// GetExtensionFromFilename extracts the extension from a filename.
// The first parameter is the filename without the extension.
// The second parameter is the extension.
func GetExtensionFromFilename(filename string) (string, string, error) {
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		return filename[:idx], filename[idx+1:], nil
	}
	return "", "", errors.New("invalid extension in name")
}

// IsValidImage checks if the provided MIME type is valid image.
func IsValidImage(mimeType string) bool {
	validMimeTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/svg+xml": true,
		"image/webp":    true,
		"image/x-icon":  true,
	}

	return validMimeTypes[strings.ToLower(mimeType)]
}

// ChunkBytes splits a byte slice into chunks of a specified size.
// The default chunk size is 4096 bytes.
func ChunkBytes(data []byte, chunkSize ...int) [][]byte {
	size := 4096
	if len(chunkSize) > 0 {
		size = chunkSize[0]
	}

	var chunks [][]byte
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}
