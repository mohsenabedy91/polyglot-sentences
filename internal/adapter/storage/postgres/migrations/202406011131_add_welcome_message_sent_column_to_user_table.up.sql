-- Add the welcome_message_sent column
ALTER TABLE users
    ADD COLUMN welcome_message_sent BOOLEAN DEFAULT FALSE;
