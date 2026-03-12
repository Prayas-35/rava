package service

import (
	"context"
	"encoding/json"

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
