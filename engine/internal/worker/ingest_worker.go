package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Prayas-35/ragkit/engine/internal/chunk"
	"github.com/Prayas-35/ragkit/engine/internal/database"
	"github.com/Prayas-35/ragkit/engine/internal/queue"
	"github.com/Prayas-35/ragkit/engine/internal/vector"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartWorker(ch *amqp.Channel, store *vector.Store) {

	msgs, err := ch.Consume(
		"rag_ingest",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for msg := range msgs {

			var job queue.IngestJob

			err := json.Unmarshal(msg.Body, &job)
			if err != nil {
				log.Println(err)
				_ = msg.Ack(false)
				continue
			}

			_ = setDocumentStatus(context.Background(), job.DocumentID, "processing")

			requeue, err := processDocument(job, store)
			if err != nil {
				log.Println("ingest error:", err, "doc=", job.DocumentID)
				if requeue {
					if nackErr := msg.Nack(false, true); nackErr != nil {
						log.Println("failed to nack message:", nackErr)
					}
					continue
				}

				_ = setDocumentStatus(context.Background(), job.DocumentID, "failed")
				if ackErr := msg.Ack(false); ackErr != nil {
					log.Println("failed to ack message:", ackErr)
				}
				continue
			}

			_ = setDocumentStatus(context.Background(), job.DocumentID, "processed")
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Println("failed to ack message:", ackErr)
			}
		}

	}()
}

func processDocument(job queue.IngestJob, store *vector.Store) (bool, error) {

	ctx := context.Background()

	parts := chunk.SplitText(job.Content, 700)
	if len(parts) == 0 {
		return false, fmt.Errorf("skipping document with empty content: %s", job.DocumentID)
	}

	vectors, err := store.Embedder.EmbedBatch(ctx, parts)
	if err != nil {
		return true, fmt.Errorf("embedding error: %w", err)
	}
	if len(vectors) != len(parts) {
		return true, fmt.Errorf("embedding error: parts and vectors length mismatch parts=%d vectors=%d doc=%s", len(parts), len(vectors), job.DocumentID)
	}

	var rows []vector.ChunkRow

	for i, text := range parts {
		if len(vectors[i]) == 0 {
			log.Println("skipping chunk with empty embedding: doc=", job.DocumentID, "chunkIndex=", i)
			continue
		}

		rows = append(rows, vector.ChunkRow{
			DocumentID: job.DocumentID,
			ProjectID:  job.ProjectID,
			Content:    text,
			Embedding:  vectors[i],
			ChunkIndex: i,
			Metadata:   job.Metadata,
		})
	}

	if len(rows) == 0 {
		return false, fmt.Errorf("no valid chunks to insert for document: %s", job.DocumentID)
	}

	err = store.BulkInsertChunks(ctx, rows)
	if err != nil {
		return true, fmt.Errorf("insert error: %w", err)
	}

	log.Println("document processed:", job.DocumentID)
	return false, nil
}

func setDocumentStatus(ctx context.Context, documentID, status string) error {
	_, err := database.DB.Exec(ctx, `UPDATE documents SET status = $1 WHERE id = $2`, status, documentID)
	if err != nil {
		log.Println("failed to update document status:", err, "doc=", documentID, "status=", status)
	}
	return err
}
