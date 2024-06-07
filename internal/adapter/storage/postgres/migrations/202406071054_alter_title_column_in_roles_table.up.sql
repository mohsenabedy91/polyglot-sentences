-- Remove the roles_title_key constraint
ALTER TABLE roles
    DROP CONSTRAINT IF EXISTS roles_title_key;