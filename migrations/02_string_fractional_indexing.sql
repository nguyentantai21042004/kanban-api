-- ============================================================================
-- STRING-BASED FRACTIONAL INDEXING SYSTEM
-- Advanced position management with automatic rebalancing
-- ============================================================================

-- ============================================================================
-- 1. REBALANCING JOB SYSTEM
-- ============================================================================

-- Job queue for rebalancing operations
CREATE TABLE IF NOT EXISTS rebalance_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID, -- NULL for list rebalancing, specific list for card rebalancing
    board_id UUID, -- For board-wide operations
    target_type VARCHAR(20) NOT NULL DEFAULT 'cards', -- 'cards' or 'lists'
    priority VARCHAR(20) DEFAULT 'background', -- critical, high, normal, low, background
    status VARCHAR(20) DEFAULT 'pending',      -- pending, running, completed, failed, cancelled
    trigger_reason VARCHAR(100),
    strategy VARCHAR(20) DEFAULT 'conservative', -- conservative, aggressive
    scheduled_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    error_message TEXT,
    result JSONB,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    
    CONSTRAINT fk_rebalance_jobs_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_rebalance_jobs_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

-- Optimized indexes for job queue processing
CREATE INDEX IF NOT EXISTS idx_rebalance_jobs_queue ON rebalance_jobs(status, priority, scheduled_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_rebalance_jobs_list ON rebalance_jobs(list_id, status);
CREATE INDEX IF NOT EXISTS idx_rebalance_jobs_board ON rebalance_jobs(board_id, status);
CREATE INDEX IF NOT EXISTS idx_rebalance_jobs_created_at ON rebalance_jobs(created_at DESC);

-- ============================================================================
-- 2. REBALANCING EVENTS & METRICS
-- ============================================================================

-- Event log for tracking rebalancing operations
CREATE TABLE IF NOT EXISTS rebalance_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID,
    board_id UUID,
    target_type VARCHAR(20) NOT NULL DEFAULT 'cards',
    strategy VARCHAR(20),
    record_count INTEGER,
    avg_length_before DECIMAL(10,4),
    avg_length_after DECIMAL(10,4),
    max_length_before INTEGER,
    max_length_after INTEGER,
    min_length_before INTEGER,
    min_length_after INTEGER,
    duration_ms INTEGER,
    trigger_reason VARCHAR(100),
    job_id UUID REFERENCES rebalance_jobs(id),
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_rebalance_events_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_rebalance_events_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_rebalance_events_list ON rebalance_events(list_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rebalance_events_board ON rebalance_events(board_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rebalance_events_job ON rebalance_events(job_id);

-- ============================================================================
-- 3. MIGRATION TRACKING
-- ============================================================================

-- Track migration progress from numeric to string positions
CREATE TABLE IF NOT EXISTS migration_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(50) NOT NULL,
    target_type VARCHAR(20) NOT NULL, -- 'cards' or 'lists'
    total_records BIGINT NOT NULL,
    migrated_records BIGINT DEFAULT 0,
    failed_records BIGINT DEFAULT 0,
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    last_updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'pending', -- pending, running, completed, failed
    error_details JSONB,
    migration_strategy VARCHAR(50) DEFAULT 'background',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_migration_progress_status ON migration_progress(status, table_name);
CREATE INDEX IF NOT EXISTS idx_migration_progress_updated ON migration_progress(last_updated_at DESC);

-- ============================================================================
-- 4. POSITION VALIDATION & MONITORING
-- ============================================================================

-- Position validation log for data integrity monitoring
CREATE TABLE IF NOT EXISTS position_validation_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID,
    board_id UUID,
    validation_type VARCHAR(50) NOT NULL, -- 'order_check', 'length_check', 'format_check', 'consistency_check'
    target_type VARCHAR(20) NOT NULL DEFAULT 'cards',
    is_valid BOOLEAN NOT NULL,
    error_message TEXT,
    error_details JSONB,
    record_count INTEGER,
    problematic_records JSONB, -- Array of IDs with issues
    checked_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_validation_log_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_validation_log_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_validation_log_list ON position_validation_log(list_id, checked_at DESC);
CREATE INDEX IF NOT EXISTS idx_validation_log_board ON position_validation_log(board_id, checked_at DESC);
CREATE INDEX IF NOT EXISTS idx_validation_log_valid ON position_validation_log(is_valid, validation_type);

-- ============================================================================
-- 5. PERFORMANCE MONITORING
-- ============================================================================

-- Position statistics for monitoring system health
CREATE TABLE IF NOT EXISTS position_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID,
    board_id UUID,
    target_type VARCHAR(20) NOT NULL DEFAULT 'cards',
    
    -- Statistics
    record_count INTEGER NOT NULL,
    avg_length DECIMAL(10,4),
    max_length INTEGER,
    min_length INTEGER,
    length_stddev DECIMAL(10,4),
    long_key_count INTEGER, -- Count of positions > 10 characters
    long_key_percentage DECIMAL(5,2),
    
    -- Health indicators
    needs_rebalance BOOLEAN DEFAULT FALSE,
    health_score DECIMAL(5,2), -- 0-100 score
    performance_impact VARCHAR(20), -- 'none', 'low', 'medium', 'high', 'critical'
    
    -- Metadata
    calculated_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '1 hour'),
    
    CONSTRAINT fk_position_stats_list FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT fk_position_stats_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_position_stats_list ON position_statistics(list_id, calculated_at DESC);
CREATE INDEX IF NOT EXISTS idx_position_stats_board ON position_statistics(board_id, calculated_at DESC);
CREATE INDEX IF NOT EXISTS idx_position_stats_health ON position_statistics(health_score, needs_rebalance);
CREATE INDEX IF NOT EXISTS idx_position_stats_expires ON position_statistics(expires_at);

-- ============================================================================
-- 6. UTILITY FUNCTIONS
-- ============================================================================

-- Function to calculate position health score
CREATE OR REPLACE FUNCTION calculate_position_health_score(
    avg_len DECIMAL,
    max_len INTEGER,
    long_key_pct DECIMAL,
    record_count INTEGER
) RETURNS DECIMAL AS $$
DECLARE
    score DECIMAL := 100.0;
BEGIN
    -- Penalize based on average length
    IF avg_len > 15 THEN
        score := score - 40;
    ELSIF avg_len > 10 THEN
        score := score - 25;
    ELSIF avg_len > 6 THEN
        score := score - 10;
    END IF;
    
    -- Penalize based on max length
    IF max_len > 25 THEN
        score := score - 30;
    ELSIF max_len > 15 THEN
        score := score - 15;
    END IF;
    
    -- Penalize based on long key percentage
    IF long_key_pct > 50 THEN
        score := score - 20;
    ELSIF long_key_pct > 25 THEN
        score := score - 10;
    END IF;
    
    -- Ensure score is not negative
    RETURN GREATEST(score, 0.0);
END;
$$ LANGUAGE plpgsql;

-- Function to determine if rebalancing is needed
CREATE OR REPLACE FUNCTION needs_rebalancing(
    avg_len DECIMAL,
    max_len INTEGER,
    long_key_pct DECIMAL
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN (
        avg_len > 12.0 OR 
        max_len > 20 OR 
        long_key_pct > 30.0
    );
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- 7. COLUMN COMMENTS FOR NEW TABLES
-- ============================================================================

-- Comments for rebalance_jobs
COMMENT ON TABLE rebalance_jobs IS 'Job queue for managing position rebalancing operations';
COMMENT ON COLUMN rebalance_jobs.target_type IS 'Type of entity to rebalance: cards or lists';
COMMENT ON COLUMN rebalance_jobs.priority IS 'Job priority: critical, high, normal, low, background';
COMMENT ON COLUMN rebalance_jobs.strategy IS 'Rebalancing strategy: conservative or aggressive';
COMMENT ON COLUMN rebalance_jobs.trigger_reason IS 'What triggered this rebalancing job';

-- Comments for rebalance_events
COMMENT ON TABLE rebalance_events IS 'Event log for tracking rebalancing operations and their results';
COMMENT ON COLUMN rebalance_events.duration_ms IS 'Time taken to complete the rebalancing operation in milliseconds';

-- Comments for migration_progress
COMMENT ON TABLE migration_progress IS 'Tracks progress of migrating from numeric to string positions';

-- Comments for position_validation_log
COMMENT ON TABLE position_validation_log IS 'Log of position validation checks for data integrity monitoring';

-- Comments for position_statistics
COMMENT ON TABLE position_statistics IS 'Cached statistics about position distribution and health';
COMMENT ON COLUMN position_statistics.health_score IS 'Position system health score (0-100)';
COMMENT ON COLUMN position_statistics.performance_impact IS 'Estimated performance impact of current position distribution';

-- ============================================================================
-- 8. INITIAL VALIDATION
-- ============================================================================

-- Validate that the system is ready for string positions
DO $$
DECLARE
    card_count INTEGER;
    list_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO card_count FROM cards WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO list_count FROM lists WHERE deleted_at IS NULL;
    
    RAISE NOTICE 'String-based fractional indexing system initialized';
    RAISE NOTICE 'Found % cards and % lists ready for migration', card_count, list_count;
    RAISE NOTICE 'System is ready for gradual migration to string positions';
END $$;
