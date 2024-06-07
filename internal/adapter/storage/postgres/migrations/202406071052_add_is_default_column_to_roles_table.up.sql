-- Add the is_default column
ALTER TABLE roles
    ADD COLUMN is_default BOOLEAN DEFAULT FALSE;