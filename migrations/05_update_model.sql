-- Rename "title" to "name" and add "alias" field to boards table
ALTER TABLE boards RENAME COLUMN title TO name;
ALTER TABLE boards ADD COLUMN alias VARCHAR(255);
