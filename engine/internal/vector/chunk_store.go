package vector

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type ChunkRow struct {
	DocumentID string
	ProjectID  string
	Content    string
	Embedding  []float32
	ChunkIndex int
	Metadata   map[string]interface{}
}

type Store struct {
	DB       *pgxpool.Pool
	Embedder *GeminiEmbedder
}

func (s *Store) BulkInsertChunks(ctx context.Context, rows []ChunkRow) error {

	var data [][]any

	for _, r := range rows {
		if len(r.Embedding) == 0 {
			log.Println("BulkInsertChunks: skipping row with empty embedding for document", r.DocumentID)
			continue
		}

		// Ensure embedding length matches the database vector dimension (1536)
		const embeddingDim = 1536
		if len(r.Embedding) != embeddingDim {
			log.Println("BulkInsertChunks: adjusting embedding length",
				"from", len(r.Embedding), "to", embeddingDim, "for document", r.DocumentID)
			vec := make([]float32, embeddingDim)
			if len(r.Embedding) > embeddingDim {
				copy(vec, r.Embedding[:embeddingDim])
			} else {
				copy(vec, r.Embedding)
			}
			r.Embedding = vec
		}

		meta, _ := json.Marshal(r.Metadata)

		vec := pgvector.NewVector(r.Embedding)

		data = append(data, []any{
			r.DocumentID,
			r.ProjectID,
			r.Content,
			vec,
			r.ChunkIndex,
			meta,
		})
	}

	if len(data) == 0 {
		log.Println("BulkInsertChunks: no rows to insert")
		return nil
	}

	_, err := s.DB.CopyFrom(
		ctx,
		pgx.Identifier{"chunks"},
		[]string{
			"document_id",
			"project_id",
			"content",
			"embedding",
			"chunk_index",
			"metadata",
		},
		pgx.CopyFromRows(data),
	)

	return err
}
