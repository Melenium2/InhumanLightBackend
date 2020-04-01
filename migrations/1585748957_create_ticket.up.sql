CREATE TABLE tickets (
    id bigserial not null PRIMARY KEY,
    title VARCHAR not null,
    description VARCHAR not null,
    section VARCHAR not null,
    from_user INT not null,
    helper INT,
    created_at TIMESTAMP,
    status varchar
);