CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" VARCHAR NOT NULL,
    "balance" BIGINT NOT NULL,
    "currency" VARCHAR NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" BIGINT NOT NULL,
    "amount" BIGINT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "reciepent_id" BIGINT NOT NULL,
    "sender_id" BIGINT NOT NULL,
    "amount" BIGINT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" 
ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;

ALTER TABLE "transfers"  
ADD FOREIGN KEY ("sender_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;

ALTER TABLE "transfers" 
ADD FOREIGN KEY ("reciepent_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;

CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("sender_id");
CREATE INDEX ON "transfers" ("reciepent_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
