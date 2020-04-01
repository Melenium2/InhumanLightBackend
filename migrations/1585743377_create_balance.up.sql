CREATE TABLE balance (
    id bigserial not null PRIMARY KEY,
    transaction_value FLOAT not null,
    balance_now FLOAT not null,
    from_market varchar,
    transaction_at TIMESTAMP,
    additional_info VARCHAR,
    user_id INTEGER
);