-- Add the email_verified_at column
ALTER TABLE users
    ADD COLUMN email_verified_at TIMESTAMP DEFAULT NULL;
