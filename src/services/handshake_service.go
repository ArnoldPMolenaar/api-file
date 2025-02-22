package services

import (
	"api-file/main/src/cache"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

// CreateHandshake creates a handshake code for the given id.
func CreateHandshake(app, id string) (string, error) {
	code, err := uuid.NewUUID()
	if err != nil {
		return "", errors.New("failed to generate handshake code")
	}
	key := handshakeCacheKey(app, code.String())

	expiration := os.Getenv("VALKEY_EXPIRATION_HANDSHAKE")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return "", err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(key).Value(id).Ex(duration).Build())
	if result.Error() != nil {
		return "", result.Error()
	}

	return code.String(), nil
}

// GetHandshake gets the id for the given app and code.
func GetHandshake(app, code string) (string, error) {
	key := handshakeCacheKey(app, code)

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return "", result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return "", err
	}

	return value, nil
}

// Creates a key for the handshake cache.
func handshakeCacheKey(app, code string) string {
	return fmt.Sprintf("handshake:%s:%s", app, code)
}
