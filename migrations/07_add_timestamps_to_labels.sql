-- Add timestamp columns to labels table
ALTER TABLE labels ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE labels ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE labels ADD COLUMN deleted_at TIMESTAMPTZ;

-- Create index for soft delete
CREATE INDEX IF NOT EXISTS idx_labels_deleted_at ON labels (deleted_at); 