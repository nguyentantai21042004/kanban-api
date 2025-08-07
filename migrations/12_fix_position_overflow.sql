-- Fix position field overflow by increasing precision
-- Change from NUMERIC(10,5) to NUMERIC(20,6) for both cards and lists

-- Update cards position field
ALTER TABLE cards ALTER COLUMN position TYPE NUMERIC(20,6);

-- Update lists position field  
ALTER TABLE lists ALTER COLUMN position TYPE NUMERIC(20,6);

-- Add comment to explain the change
COMMENT ON COLUMN cards.position IS 'Card position using fractional indexing - supports large values up to 99999999999999.999999';
COMMENT ON COLUMN lists.position IS 'List position using fractional indexing - supports large values up to 99999999999999.999999';g