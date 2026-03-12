package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Prayas-35/ragkit/engine/internal/database"
)

type APIKey struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
}

// generateRawKey creates a developer-facing API key like "rag_...".
func generateRawKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("rag_%s", hex.EncodeToString(buf)), nil
}

func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// CreateAPIKey creates a new API key for a project and returns the raw key once.
func CreateAPIKey(ctx context.Context, projectID string) (string, *APIKey, error) {
	raw, err := generateRawKey()
	if err != nil {
		return "", nil, err
	}

	keyHash := hashKey(raw)
	row := database.DB.QueryRow(ctx,
		`INSERT INTO api_keys (project_id, key_hash) VALUES ($1, $2)
		RETURNING id, project_id`,
		projectID,
		keyHash,
	)

	var k APIKey
	if err := row.Scan(&k.ID, &k.ProjectID); err != nil {
		return "", nil, err
	}

	return raw, &k, nil
}

// ResolveProjectIDByAPIKey finds the project_id associated with a raw API key, if valid and not revoked.
func ResolveProjectIDByAPIKey(ctx context.Context, rawKey string) (string, error) {
	if strings.TrimSpace(rawKey) == "" {
		return "", fmt.Errorf("api key is required")
	}

	hash := hashKey(rawKey)
	row := database.DB.QueryRow(ctx,
		`SELECT project_id FROM api_keys WHERE key_hash = $1 AND revoked = FALSE`,
		hash,
	)

	var projectID string
	if err := row.Scan(&projectID); err != nil {
		return "", err
	}

	return projectID, nil
}
