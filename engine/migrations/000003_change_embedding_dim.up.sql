-- No-op migration: keep embedding dimension at 1536 for HNSW index compatibility
-- HNSW in pgvector supports up to 2000 dimensions, and the current
-- schema already uses vector(1536), which is valid.

SELECT 1;
