-- Add content_hash column for idempotent ingestion
-- This allows us to detect when identical content is being ingested and return the existing document ID

ALTER TABLE documents ADD COLUMN content_hash TEXT;

-- Create a unique index on (project_id, content_hash) to enforce one document per unique content per project
CREATE UNIQUE INDEX documents_project_id_content_hash_idx 
ON documents (project_id, content_hash) 
WHERE content_hash IS NOT NULL;
