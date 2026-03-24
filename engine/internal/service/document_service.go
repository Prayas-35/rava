package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/Prayas-35/ragkit/engine/internal/database"
	"github.com/Prayas-35/ragkit/engine/internal/queue"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DocumentIngestRequest struct {
	Name     string                 `json:"name"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
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

	job := queue.IngestJob{
		DocumentID: docID,
		ProjectID:  projectID,
		Content:    req.Content,
		Metadata:   req.Metadata,
	}

	body, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	// Publish job to RabbitMQ
	if err := ch.Publish(
		"",
		"rag_ingest",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return "", err
	}

	return docID, nil
}

// calculateContentHash computes SHA256 hash of content for idempotent ingestion detection
func calculateContentHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

// IngestDocumentIdempotent creates a document record with idempotent upsert logic.
// If a document with the same (projectID, contentHash) already exists, it returns that document ID.
// Otherwise, it creates a new document and enqueues the ingestion job.
func IngestDocumentIdempotent(ch *amqp.Channel, projectID string, req DocumentIngestRequest) (string, error) {
	ctx := context.Background()

	// Calculate content hash for idempotency detection
	contentHash := calculateContentHash(req.Content)

	// Check if a document with this content already exists for this project
	var existingDocID string
	err := database.DB.QueryRow(ctx,
		`SELECT id FROM documents WHERE project_id = $1 AND content_hash = $2`,
		projectID,
		contentHash,
	).Scan(&existingDocID)

	if err == nil {
		// Document with this content already exists - return existing ID (idempotent)
		return existingDocID, nil
	}

	// Document doesn't exist, create a new one
	docID := uuid.New().String()

	// Insert document row with content_hash so that chunks can reference it via FK
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

	job := queue.IngestJob{
		DocumentID: docID,
		ProjectID:  projectID,
		Content:    req.Content,
		Metadata:   req.Metadata,
	}

	body, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	// Publish job to RabbitMQ
	if err := ch.Publish(
		"",
		"rag_ingest",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return "", err
	}

	return docID, nil
}
