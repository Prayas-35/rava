-- enable required extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

--------------------------------------------------
-- USERS
--------------------------------------------------

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL UNIQUE,
    name        TEXT,
    password_hash TEXT,
    created_at  TIMESTAMPTZ DEFAULT now()
);

--------------------------------------------------
-- PROJECTS
--------------------------------------------------

CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now()
);

--------------------------------------------------
-- API KEYS
--------------------------------------------------

CREATE TABLE api_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key_hash    TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ DEFAULT now(),
    revoked     BOOLEAN DEFAULT FALSE
);

--------------------------------------------------
-- DOCUMENTS
--------------------------------------------------

CREATE TABLE documents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending',
    metadata    JSONB,
    created_at  TIMESTAMPTZ DEFAULT now()
);

--------------------------------------------------
-- CHUNKS
--------------------------------------------------

CREATE TABLE chunks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id  UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    content      TEXT NOT NULL,
    embedding    vector(1536),
    chunk_index  INT NOT NULL,
    metadata     JSONB,
    created_at   TIMESTAMPTZ DEFAULT now()
);

--------------------------------------------------
-- VECTOR INDEX FOR SEARCH
--------------------------------------------------

CREATE INDEX chunks_embedding_idx
ON chunks
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
