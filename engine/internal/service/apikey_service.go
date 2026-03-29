package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/database"
)

type APIKey struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
}

type APIKeyRecord struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
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

// CreateAPIKey creates a new API key for a project and stores an encrypted copy for future retrieval.
func CreateAPIKey(ctx context.Context, projectID, encryptionSecret string) (string, *APIKey, error) {
	if strings.TrimSpace(encryptionSecret) == "" {
		return "", nil, errors.New("API key encryption secret is required")
	}

	raw, err := generateRawKey()
	if err != nil {
		return "", nil, err
	}

	keyHash := hashKey(raw)
	row := database.DB.QueryRow(ctx,
		`INSERT INTO api_keys (project_id, key_hash, key_encrypted) VALUES ($1, $2, pgp_sym_encrypt($3, $4))
		RETURNING id, project_id`,
		projectID,
		keyHash,
		raw,
		encryptionSecret,
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

func ListAPIKeysByProject(ctx context.Context, projectID, encryptionSecret string) ([]APIKeyRecord, error) {
	if strings.TrimSpace(encryptionSecret) == "" {
		return nil, errors.New("API key encryption secret is required")
	}

	rows, err := database.DB.Query(ctx,
		`SELECT id, project_id, pgp_sym_decrypt(key_encrypted, $2)::text as api_key, created_at, revoked
		 FROM api_keys
		 WHERE project_id = $1 AND key_encrypted IS NOT NULL
		 ORDER BY created_at DESC`,
		projectID,
		encryptionSecret,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]APIKeyRecord, 0)
	for rows.Next() {
		var key APIKeyRecord
		if err := rows.Scan(&key.ID, &key.ProjectID, &key.APIKey, &key.CreatedAt, &key.Revoked); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
