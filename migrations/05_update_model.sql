-- Rename "title" to "name" and add "alias" field to boards table
ALTER TABLE boards RENAME COLUMN title TO name;
ALTER TABLE boards ADD COLUMN alias VARCHAR(255);

-- Add created_by, updated_by, deleted_by fields to labels table
ALTER TABLE labels ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE labels ADD COLUMN updated_by UUID REFERENCES users(id);
ALTER TABLE labels ADD COLUMN deleted_by UUID REFERENCES users(id);

ALTER TABLE users ADD COLUMN role_id UUID REFERENCES roles(id);

ALTER TABLE boards ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE lists ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE cards ADD COLUMN created_by UUID REFERENCES users(id);
