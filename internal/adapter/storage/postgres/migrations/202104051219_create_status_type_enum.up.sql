-- Create a new type called status_type that is an enumeration
-- of the values 'ACTIVE', 'DISABLED', 'UNPUBLISHED' and 'DRAFT'.
CREATE TYPE status_type AS ENUM ('ACTIVE', 'DISABLED', 'UNPUBLISHED', 'DRAFT');