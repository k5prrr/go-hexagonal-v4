CREATE TABLE users (
    uid TEXT PRIMARY KEY,
    family_name TEXT NOT NULL,
    name TEXT NOT NULL,
    middle_name TEXT,
    description TEXT,
    birth_date DATE,
    phone TEXT,
    email TEXT,
    phone_confirmed BOOLEAN DEFAULT FALSE,
    email_confirmed BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    password_hash TEXT NOT NULL,
    key_api TEXT NOT NULL
);
