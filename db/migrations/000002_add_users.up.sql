CREATE TABLE users (
    username VARCHAR PRIMARY KEY,
    hashed_password VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_changed_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users(username);

ALTER TABLE accounts ADD CONSTRAINT "owner_currency_key" UNIQUE (owner, currency);