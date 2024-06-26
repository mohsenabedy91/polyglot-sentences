-- Table: grammars
CREATE TABLE IF NOT EXISTS grammars
(
    id         INTEGER GENERATED BY DEFAULT AS IDENTITY
        CONSTRAINT pk_grammars PRIMARY KEY,
    uuid       uuid                     DEFAULT gen_random_uuid() UNIQUE,
    title      VARCHAR(255) NOT NULL UNIQUE,
    status     status_type              DEFAULT 'ACTIVE'::status_type,
    created_by INTEGER
        CONSTRAINT fk_grammars_created_by REFERENCES users,
    updated_by INTEGER
        CONSTRAINT fk_grammars_updated_by REFERENCES users,
    deleted_by INTEGER
        CONSTRAINT fk_grammars_deleted_by REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);