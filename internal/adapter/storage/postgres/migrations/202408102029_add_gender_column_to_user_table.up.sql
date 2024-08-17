-- Add the gender column
ALTER TABLE users
    ADD COLUMN gender user_gender_type DEFAULT 'PREFER_NOT_TO_SAY'::user_gender_type;