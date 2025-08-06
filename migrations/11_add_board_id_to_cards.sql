-- Add board_id field to cards table
-- This field will provide direct reference to the board, similar to list_id

-- Step 1: Add board_id column as nullable first
ALTER TABLE cards ADD COLUMN board_id UUID;

-- Step 2: Populate board_id with data from existing cards by joining through lists table
UPDATE cards SET board_id = (
    SELECT l.board_id 
    FROM lists l 
    WHERE l.id = cards.list_id
);

-- Step 3: Now make it NOT NULL (like list_id)
ALTER TABLE cards ALTER COLUMN board_id SET NOT NULL;

-- Step 4: Add foreign key constraint (ON DELETE CASCADE like list_id)
ALTER TABLE cards ADD CONSTRAINT fk_cards_board 
    FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE;

-- Step 5: Create indexes for better query performance (similar to list_id indexes)
CREATE INDEX IF NOT EXISTS idx_cards_board_position ON cards (board_id, position);
CREATE INDEX IF NOT EXISTS idx_cards_board_active ON cards (board_id, is_archived);
CREATE INDEX IF NOT EXISTS idx_cards_board_priority ON cards (board_id, priority) WHERE is_archived = FALSE;

-- Add comment to explain the field
COMMENT ON COLUMN cards.board_id IS 'Board this card belongs to - provides direct reference without joining through lists';