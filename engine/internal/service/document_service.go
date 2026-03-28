package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Prayas-35/ragkit/engine/internal/database"
	"github.com/Prayas-35/ragkit/engine/internal/queue"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DocumentIngestRequest struct {
	Name     string                 `json:"name"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

func enqueueIngestionJob(ch *amqp.Channel, docID, projectID string, req DocumentIngestRequest) error {
	job := queue.IngestJob{
		DocumentID: docID,
		ProjectID:  projectID,
		Content:    req.Content,
		Metadata:   req.Metadata,
	}

	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		"rag_ingest",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// QueueDocumentIngestion creates a document record and enqueues a job for ingestion.
func QueueDocumentIngestion(ch *amqp.Channel, projectID string, req DocumentIngestRequest) (string, error) {
	ctx := context.Background()

	// Generate a new document ID
	docID := uuid.New().String()

	// Insert document row so that chunks can reference it via FK
	_, err := database.DB.Exec(ctx,
		`INSERT INTO documents (id, project_id, name, status, metadata) VALUES ($1, $2, $3, $4, $5)`,
		docID,
		projectID,
		req.Name,
		"pending",
		req.Metadata,
	)
	if err != nil {
		return "", err
	}

	if err := enqueueIngestionJob(ch, docID, projectID, req); err != nil {
		return "", err
	}

	return docID, nil
}

// calculateContentHash computes SHA256 hash of content for idempotent ingestion detection
func calculateContentHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

// IngestDocumentIdempotent upserts ingestion by (projectID, name).
// If a document with the same name already exists in the project, it updates that row
// and reuses the same document ID. Otherwise, it creates a new row.
func IngestDocumentIdempotent(ch *amqp.Channel, projectID string, req DocumentIngestRequest) (string, error) {
	ctx := context.Background()

	// Calculate content hash for tracking/debugging.
	contentHash := calculateContentHash(req.Content)

	// Reuse the existing document row for this (projectID, name).
	var existingDocID string
	err := database.DB.QueryRow(ctx,
		`SELECT id FROM documents WHERE project_id = $1 AND name = $2 ORDER BY created_at DESC LIMIT 1`,
		projectID,
		req.Name,
	).Scan(&existingDocID)

	if err == nil {
		// Update existing row and clear previous chunks so re-ingestion replaces content.
		_, err = database.DB.Exec(ctx,
			`UPDATE documents SET status = $1, metadata = $2, content_hash = $3 WHERE id = $4`,
			"pending",
			req.Metadata,
			contentHash,
			existingDocID,
		)
		if err != nil {
			return "", err
		}

		_, err = database.DB.Exec(ctx, `DELETE FROM chunks WHERE document_id = $1`, existingDocID)
		if err != nil {
			return "", err
		}

		if err := enqueueIngestionJob(ch, existingDocID, projectID, req); err != nil {
			return "", err
		}

		return existingDocID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	// Document doesn't exist, create a new one.
	docID := uuid.New().String()

	// Insert document row with content_hash so that chunks can reference it via FK.
	_, err = database.DB.Exec(ctx,
		`INSERT INTO documents (id, project_id, name, status, metadata, content_hash) VALUES ($1, $2, $3, $4, $5, $6)`,
		docID,
		projectID,
		req.Name,
		"pending",
		req.Metadata,
		contentHash,
	)
	if err != nil {
		return "", err
	}

	if err := enqueueIngestionJob(ch, docID, projectID, req); err != nil {
		return "", err
	}

	return docID, nil
}
