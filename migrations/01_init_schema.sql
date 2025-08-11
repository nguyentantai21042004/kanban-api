-- ============================================================================
-- KANBAN API - COMPLETE SCHEMA INITIALIZATION
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- 1. CORE USER MANAGEMENT
-- ============================================================================

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    alias VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(100),
    password_hash VARCHAR(255),
    avatar_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    role_id UUID REFERENCES roles(id),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- ============================================================================
-- 2. KANBAN CORE ENTITIES
-- ============================================================================

-- Boards table
CREATE TABLE IF NOT EXISTS boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255),
    description TEXT,
    created_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Lists table (columns in board)
CREATE TABLE IF NOT EXISTS lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    position NUMERIC(20,6) NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_lists_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

ALTER TABLE lists
    ALTER COLUMN position TYPE VARCHAR(255);

-- Create ENUM types
CREATE TYPE card_priority AS ENUM ('low', 'medium', 'high');
CREATE TYPE card_action_type AS ENUM ('created', 'moved', 'updated', 'commented');

-- Cards table with all enhanced fields
CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL,
    board_id UUID NOT NULL,
    name VARCHAR(500) NOT NULL,
    alias VARCHAR(100),
    description TEXT,
    
    -- Position and ordering
    position NUMERIC(20,6) NOT NULL,
    
    -- Dates and scheduling
    due_date TIMESTAMPTZ,
    start_date TIMESTAMPTZ,
    completion_date TIMESTAMPTZ,
    last_activity_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    -- Categorization and priority
    priority card_priority NOT NULL DEFAULT 'medium',
    labels JSONB, -- Array of label objects
    tags TEXT[], -- Array of text tags for flexible categorization
    
    -- Assignment and tracking
    assigned_to UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    
    -- Time tracking
    estimated_hours NUMERIC(5,2),
    actual_hours NUMERIC(5,2),
    
    -- Rich content
    attachments JSONB DEFAULT '[]'::jsonb, -- JSON array of uploaded file UUIDs
    checklist JSONB DEFAULT '[]'::jsonb, -- JSON array of checklist items
    
    -- Status
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_cards_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_cards_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

-- Change position column from NUMERIC(20,6) to VARCHAR(32) to store string-based fractional index
ALTER TABLE cards
    ALTER COLUMN position TYPE VARCHAR(255);


-- ============================================================================
-- 3. SUPPORTING ENTITIES
-- ============================================================================

-- Labels master table
CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) NOT NULL,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_by UUID REFERENCES users(id),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_labels_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE,
    CONSTRAINT unique_board_label UNIQUE (board_id, name)
);

-- Card activities/audit trail
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

-- Comments table for card comments
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    parent_id UUID, -- For reply comments
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    edited_by UUID REFERENCES users(id),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_comments_card FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

-- Uploads table for file management
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
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT fk_uploads_created_user FOREIGN KEY (created_user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- ============================================================================
-- 4. PERFORMANCE INDEXES
-- ============================================================================

-- Lists indexes
CREATE INDEX IF NOT EXISTS idx_lists_board_position ON lists (board_id, position);
CREATE INDEX IF NOT EXISTS idx_lists_board_active ON lists (board_id, is_archived);

-- Cards indexes - Core functionality
CREATE INDEX IF NOT EXISTS idx_cards_list_position ON cards (list_id, position);
CREATE INDEX IF NOT EXISTS idx_cards_list_active ON cards (list_id, is_archived);
CREATE INDEX IF NOT EXISTS idx_cards_board_position ON cards (board_id, position);
CREATE INDEX IF NOT EXISTS idx_cards_board_active ON cards (board_id, is_archived);

-- Cards indexes - Enhanced features
CREATE INDEX IF NOT EXISTS idx_cards_due_date ON cards (due_date);
CREATE INDEX IF NOT EXISTS idx_cards_due_date_active ON cards (due_date) WHERE is_archived = FALSE;
CREATE INDEX IF NOT EXISTS idx_cards_priority_active ON cards (priority) WHERE is_archived = FALSE;
CREATE INDEX IF NOT EXISTS idx_cards_assigned_to ON cards (assigned_to);
CREATE INDEX IF NOT EXISTS idx_cards_last_activity ON cards (last_activity_at DESC);
CREATE INDEX IF NOT EXISTS idx_cards_completion_date ON cards (completion_date);
CREATE INDEX IF NOT EXISTS idx_cards_alias ON cards (alias) WHERE alias IS NOT NULL;

-- Activities indexes
CREATE INDEX IF NOT EXISTS idx_activities_card_time ON card_activities (card_id, created_at);

-- Comments indexes
CREATE INDEX IF NOT EXISTS idx_comments_card_id ON comments (card_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

-- Labels indexes
CREATE INDEX IF NOT EXISTS idx_labels_deleted_at ON labels (deleted_at);

-- Uploads indexes
CREATE INDEX IF NOT EXISTS idx_uploads_created_user_id ON uploads(created_user_id);
CREATE INDEX IF NOT EXISTS idx_uploads_public_id ON uploads(public_id);
CREATE INDEX IF NOT EXISTS idx_uploads_deleted_at ON uploads(deleted_at);
CREATE INDEX IF NOT EXISTS idx_uploads_bucket_object ON uploads(bucket_name, object_name);

-- ============================================================================
-- 5. COLUMN COMMENTS
-- ============================================================================

-- Cards table comments
COMMENT ON COLUMN cards.position IS 'Card position using fractional indexing - supports large values up to 99999999999999.999999';
COMMENT ON COLUMN cards.assigned_to IS 'User ID of the person assigned to this card';
COMMENT ON COLUMN cards.attachments IS 'JSON array of uploaded file UUIDs';
COMMENT ON COLUMN cards.estimated_hours IS 'Estimated time to complete the card in hours';
COMMENT ON COLUMN cards.actual_hours IS 'Actual time spent on the card in hours';
COMMENT ON COLUMN cards.start_date IS 'When work on this card should start';
COMMENT ON COLUMN cards.completion_date IS 'When the card was actually completed';
COMMENT ON COLUMN cards.tags IS 'Array of text tags for flexible categorization';
COMMENT ON COLUMN cards.checklist IS 'JSON array of checklist items with completion status';
COMMENT ON COLUMN cards.last_activity_at IS 'Timestamp of last activity on this card';
COMMENT ON COLUMN cards.updated_by IS 'User ID who last modified this card';
COMMENT ON COLUMN cards.alias IS 'Short, memorable identifier for the card (e.g., PROJ-123, BUG-001)';
COMMENT ON COLUMN cards.board_id IS 'Board this card belongs to - provides direct reference without joining through lists';

-- Lists table comments
COMMENT ON COLUMN lists.position IS 'List position using fractional indexing - supports large values up to 99999999999999.999999';

-- Comments table comments
COMMENT ON TABLE comments IS 'Comments for cards with support for threaded replies';
COMMENT ON COLUMN comments.card_id IS 'Card this comment belongs to';
COMMENT ON COLUMN comments.user_id IS 'User who created this comment';
COMMENT ON COLUMN comments.content IS 'Comment content';
COMMENT ON COLUMN comments.parent_id IS 'Parent comment ID for replies';
COMMENT ON COLUMN comments.is_edited IS 'Whether this comment has been edited';
COMMENT ON COLUMN comments.edited_at IS 'When this comment was last edited';
COMMENT ON COLUMN comments.edited_by IS 'User who last edited this comment';

-- ============================================================================
-- 6. INITIAL DATA
-- ============================================================================

-- Insert initial roles
INSERT INTO roles (name, code, alias, description) VALUES
    ('Super Admin', 'SUPER_ADMIN', 'super_admin', 'System administrator with full access'),
    ('User', 'USER', 'user', 'Regular user with basic access')
ON CONFLICT (code) DO NOTHING;
