-- Table: languages
CREATE TABLE IF NOT EXISTS languages
(
    id         INTEGER GENERATED BY DEFAULT AS IDENTITY
        CONSTRAINT pk_languages PRIMARY KEY,
    uuid       uuid                     DEFAULT gen_random_uuid() UNIQUE,
    name       VARCHAR(64) NOT NULL,
    code       VARCHAR(4)  NOT NULL,
    status     status_type              DEFAULT 'ACTIVE'::status_type,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Index: idx_languages_code
CREATE INDEX IF NOT EXISTS idx_languages_name
    ON languages (name);

-- Index: idx_languages_code
CREATE INDEX IF NOT EXISTS idx_languages_code
    ON languages (code);