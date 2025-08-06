-- Boards table
CREATE TABLE IF NOT EXISTS boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    Name VARCHAR(255) NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Lists table (columns in board)
CREATE TABLE IF NOT EXISTS lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL,
    Name VARCHAR(255) NOT NULL,
    position NUMERIC(10,5) NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_lists_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_lists_board_position ON lists (board_id, position);
CREATE INDEX IF NOT EXISTS idx_lists_board_active ON lists (board_id, is_archived);

-- Cards with rich metadata
-- Use TEXT for labels if JSON is not supported, or JSONB for better performance
CREATE TYPE card_priority AS ENUM ('low', 'medium', 'high');
CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL,
    Name VARCHAR(500) NOT NULL,
    description TEXT,
    position NUMERIC(10,5) NOT NULL,
    due_date TIMESTAMPTZ,
    priority card_priority NOT NULL DEFAULT 'medium',
    labels JSONB, -- Array of label objects
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_cards_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cards_list_position ON cards (list_id, position);
CREATE INDEX IF NOT EXISTS idx_cards_list_active ON cards (list_id, is_archived);
CREATE INDEX IF NOT EXISTS idx_cards_due_date ON cards (due_date);

-- Card activities/audit trail
CREATE TYPE card_action_type AS ENUM ('created', 'moved', 'updated', 'commented');
CREATE TABLE IF NOT EXISTS card_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL,
    action_type card_action_type NOT NULL,
    old_data JSONB,
    new_data JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_activities_card FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_activities_card_time ON card_activities (card_id, created_at);

-- Labels master table
CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) NOT NULL,
    
    CONSTRAINT fk_labels_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE,
    CONSTRAINT unique_board_label UNIQUE (board_id, name)
);