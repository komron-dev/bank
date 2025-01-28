CREATE TABLE verify_emails (
                                 id            BIGSERIAL PRIMARY KEY,
                                 username      VARCHAR NOT NULL,
                                 email         VARCHAR NOT NULL,
                                 secret_code   VARCHAR NOT NULL,
                                 is_used       BOOLEAN NOT NULL DEFAULT FALSE,
                                 created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                 expired_at    TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '15 minutes')
);

ALTER TABLE verify_emails
    ADD FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE users
    ADD COLUMN is_email_verified BOOLEAN NOT NULL DEFAULT FALSE;