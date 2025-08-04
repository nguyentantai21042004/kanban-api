-- Update cards table with new fields for better functionality
-- Add assigned_to field for card assignment
ALTER TABLE cards ADD COLUMN assigned_to UUID REFERENCES users(id);

-- Add attachments field to store list of uploaded file UUIDs
ALTER TABLE cards ADD COLUMN attachments JSONB DEFAULT '[]'::jsonb;

-- Add estimated_hours for time tracking
ALTER TABLE cards ADD COLUMN estimated_hours NUMERIC(5,2);

-- Add actual_hours for time tracking
ALTER TABLE cards ADD COLUMN actual_hours NUMERIC(5,2);

-- Add start_date for project planning
ALTER TABLE cards ADD COLUMN start_date TIMESTAMPTZ;

-- Add completion_date for tracking when card was completed
ALTER TABLE cards ADD COLUMN completion_date TIMESTAMPTZ;

-- Add tags field for flexible categorization (alternative to labels)
ALTER TABLE cards ADD COLUMN tags TEXT[];

-- Add checklist field for task breakdown
ALTER TABLE cards ADD COLUMN checklist JSONB DEFAULT '[]'::jsonb;

-- Add last_activity_at for tracking recent activity
ALTER TABLE cards ADD COLUMN last_activity_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- Add updated_by field to track who last modified the card
ALTER TABLE cards ADD COLUMN updated_by UUID REFERENCES users(id);

-- Create comments table for card comments
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

-- Create indexes for comments table
CREATE INDEX IF NOT EXISTS idx_comments_card_id ON comments (card_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_cards_assigned_to ON cards (assigned_to);
CREATE INDEX IF NOT EXISTS idx_cards_due_date_active ON cards (due_date) WHERE is_archived = FALSE;
CREATE INDEX IF NOT EXISTS idx_cards_priority_active ON cards (priority) WHERE is_archived = FALSE;
CREATE INDEX IF NOT EXISTS idx_cards_last_activity ON cards (last_activity_at DESC);
CREATE INDEX IF NOT EXISTS idx_cards_completion_date ON cards (completion_date);

-- Add comments to explain the new fields
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

-- Add comments for comments table
COMMENT ON TABLE comments IS 'Comments for cards with support for threaded replies';
COMMENT ON COLUMN comments.card_id IS 'Card this comment belongs to';
COMMENT ON COLUMN comments.user_id IS 'User who created this comment';
COMMENT ON COLUMN comments.content IS 'Comment content';
COMMENT ON COLUMN comments.parent_id IS 'Parent comment ID for replies';
COMMENT ON COLUMN comments.is_edited IS 'Whether this comment has been edited';
COMMENT ON COLUMN comments.edited_at IS 'When this comment was last edited';
COMMENT ON COLUMN comments.edited_by IS 'User who last edited this comment'; 