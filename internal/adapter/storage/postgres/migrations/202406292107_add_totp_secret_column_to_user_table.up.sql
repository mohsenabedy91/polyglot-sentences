-- Add the totp_secret column
ALTER TABLE users
    ADD COLUMN totp_secret VARCHAR(255) UNIQUE DEFAULT NULL;