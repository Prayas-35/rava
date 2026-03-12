-- Switch vector index from IVF_FLAT to HNSW

DROP INDEX IF EXISTS chunks_embedding_idx;

CREATE INDEX chunks_embedding_idx
ON chunks
USING hnsw (embedding vector_cosine_ops);
