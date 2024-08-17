-- Create a new type called user_gender_type that is an enumeration
-- of the values 'MALE', 'FEMALE', 'OTHER' and 'PREFER_NOT_TO_SAY'.
CREATE TYPE user_gender_type AS ENUM ('MALE', 'FEMALE', 'OTHER', 'PREFER_NOT_TO_SAY');