-- Migration to rename all Name columns to name
-- This migration renames Name to name for all relevant tables

-- Rename Name to name in lists table
ALTER TABLE lists RENAME COLUMN Name TO name;

-- Rename Name to name in cards table  
ALTER TABLE cards RENAME COLUMN Name TO name;

-- Update any JSON data in card_activities table that references Name
UPDATE card_activities 
SET new_data = jsonb_set(new_data, '{name}', new_data->'Name')
WHERE new_data ? 'Name';

UPDATE card_activities 
SET new_data = new_data - 'Name'
WHERE new_data ? 'Name';

UPDATE card_activities 
SET old_data = jsonb_set(old_data, '{name}', old_data->'Name')
WHERE old_data ? 'Name';

UPDATE card_activities 
SET old_data = old_data - 'Name'
WHERE old_data ? 'Name';