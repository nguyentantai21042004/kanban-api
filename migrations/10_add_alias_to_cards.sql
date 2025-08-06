-- Add alias field to cards table
-- This field will provide a short, memorable identifier for cards

ALTER TABLE cards ADD COLUMN alias VARCHAR(100);

-- Create unique index for alias within a board context
-- We need to ensure alias is unique within a board, not globally
-- First, we need to create an index that includes board context through list relationship
CREATE INDEX IF NOT EXISTS idx_cards_alias ON cards (alias) WHERE alias IS NOT NULL;

-- Add comment to explain the field
COMMENT ON COLUMN cards.alias IS 'Short, memorable identifier for the card (e.g., PROJ-123, BUG-001)';

-- Note: If you want to enforce uniqueness per board, you would need a more complex constraint
-- that joins through the lists table to get the board_id. For now, we'll allow duplicate aliases
-- across different boards but recommend making them unique within each board at the application level.