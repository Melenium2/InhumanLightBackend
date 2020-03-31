CREATE TABLE users(
    id bigserial not null PRIMARY KEY,
    username VARCHAR not null,
    email VARCHAR not null UNIQUE,
    encrypted_password varchar not null,
    created_at TIMESTAMP,
    token VARCHAR not null,
    contacts VARCHAR,
    role VARCHAR not null,
    is_active boolean
);