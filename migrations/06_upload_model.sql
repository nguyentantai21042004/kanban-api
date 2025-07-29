-- Optimized for MinIO: store MinIO-specific object info, avoid storing file path, add bucket/object/etag, and relevant metadata

CREATE TABLE IF NOT EXISTS uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket_name VARCHAR(100) NOT NULL,
    object_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    etag VARCHAR(255),
    metadata JSONB,
    url TEXT,
    source VARCHAR(100) NOT NULL,
    public_id VARCHAR(255),
    created_user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_uploads_created_user_id ON uploads(created_user_id);
CREATE INDEX IF NOT EXISTS idx_uploads_public_id ON uploads(public_id);
CREATE INDEX IF NOT EXISTS idx_uploads_deleted_at ON uploads(deleted_at);
CREATE INDEX IF NOT EXISTS idx_uploads_bucket_object ON uploads(bucket_name, object_name);

ALTER TABLE uploads
    ADD CONSTRAINT fk_uploads_created_user
    FOREIGN KEY (created_user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;