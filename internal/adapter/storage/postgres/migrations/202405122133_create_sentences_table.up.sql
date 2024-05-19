-- Table: sentences
CREATE TABLE IF NOT EXISTS sentences
(
    id         INTEGER GENERATED BY DEFAULT AS IDENTITY
        CONSTRAINT pk_sentences PRIMARY KEY,
    uuid       uuid                     DEFAULT gen_random_uuid() UNIQUE,
    TEXT       TEXT    NOT NULL,
    status     status_type              DEFAULT 'ACTIVE'::status_type,
    level      sentence_level_type      DEFAULT 'EASY'::sentence_level_type,
    grammar_id INTEGER NOT NULL,
    created_by INTEGER
        CONSTRAINT fk_sentences_created_by REFERENCES users,
    updated_by INTEGER
        CONSTRAINT fk_sentences_updated_by REFERENCES users,
    deleted_by INTEGER
        CONSTRAINT fk_sentences_deleted_by REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (grammar_id) REFERENCES grammars (id)
)