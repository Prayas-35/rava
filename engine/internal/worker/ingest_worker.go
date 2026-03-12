package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Prayas-35/ragkit/engine/internal/chunk"
	"github.com/Prayas-35/ragkit/engine/internal/queue"
	"github.com/Prayas-35/ragkit/engine/internal/vector"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartWorker(ch *amqp.Channel, store *vector.Store) {

	msgs, err := ch.Consume(
		"rag_ingest",
		"",
		true,
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
				continue
			}

			processDocument(job, store)
		}

	}()
}

func processDocument(job queue.IngestJob, store *vector.Store) {

	ctx := context.Background()

	parts := chunk.SplitText(job.Content, 700)
	if len(parts) == 0 {
		log.Println("skipping document with empty content:", job.DocumentID)
		return
	}

	vectors, err := store.Embedder.EmbedBatch(ctx, parts)
	if err != nil {
		log.Println("embedding error:", err)
		return
	}
	if len(vectors) != len(parts) {
		log.Println("embedding error: parts and vectors length mismatch",
			"parts=", len(parts), "vectors=", len(vectors), "doc=", job.DocumentID)
		return
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
		log.Println("no valid chunks to insert for document:", job.DocumentID)
		return
	}

	err = store.BulkInsertChunks(ctx, rows)
	if err != nil {
		log.Println("insert error:", err)
		return
	}

	log.Println("document processed:", job.DocumentID)
}
