-- Create a new type called user_status_type that is an enumeration
-- of the values 'ACTIVE', 'INACTIVE', 'UNVERIFIED' and 'BANNED'.
CREATE TYPE user_status_type AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'UNVERIFIED',
    'BANNED'
    );
