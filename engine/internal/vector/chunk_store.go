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

type SearchResult struct {
	Content    string
	DocumentID string
	ChunkIndex int
	Score      float32
}

const embeddingDim = 1536

func normalizeEmbeddingDim(embedding []float32) []float32 {
	if len(embedding) == embeddingDim {
		return embedding
	}

	vec := make([]float32, embeddingDim)
	if len(embedding) > embeddingDim {
		copy(vec, embedding[:embeddingDim])
		return vec
	}

	copy(vec, embedding)
	return vec
}

func (s *Store) BulkInsertChunks(ctx context.Context, rows []ChunkRow) error {

	var data [][]any

	for _, r := range rows {
		if len(r.Embedding) == 0 {
			log.Println("BulkInsertChunks: skipping row with empty embedding for document", r.DocumentID)
			continue
		}

		if len(r.Embedding) != embeddingDim {
			log.Println("BulkInsertChunks: adjusting embedding length",
				"from", len(r.Embedding), "to", embeddingDim, "for document", r.DocumentID)
			r.Embedding = normalizeEmbeddingDim(r.Embedding)
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

func (s *Store) VectorSearch(
	ctx context.Context,
	projectID string,
	query string,
	limit int,
) ([]SearchResult, error) {

	// 1. Embed query
	embedding, err := s.Embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(embedding) != embeddingDim {
		embedding = normalizeEmbeddingDim(embedding)
	}

	vec := pgvector.NewVector(embedding)

	// 2. Run vector similarity search
	rows, err := s.DB.Query(
		ctx,
		`
		SELECT 
			content,
			document_id,
			chunk_index,
			embedding <-> $1 AS score
		FROM chunks
		WHERE project_id = $2
		ORDER BY embedding <-> $1
		LIMIT $3
		`,
		vec,
		projectID,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []SearchResult

	for rows.Next() {

		var r SearchResult

		err := rows.Scan(
			&r.Content,
			&r.DocumentID,
			&r.ChunkIndex,
			&r.Score,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if len(results) > 0 {
		return results, nil
	}

	rows, err = s.DB.Query(
		ctx,
		`
		SELECT
			content,
			document_id,
			chunk_index,
			0.0 AS score
		FROM chunks
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		`,
		projectID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r SearchResult
		err := rows.Scan(
			&r.Content,
			&r.DocumentID,
			&r.ChunkIndex,
			&r.Score,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}
