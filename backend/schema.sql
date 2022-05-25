CREATE TABLE "accounts" (
    "username" varchar PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "hashed_password" bytea NOT NULL,
    "is_deleted" boolean NOT NULL DEFAULT false,
    "is_email_verified" boolean NOT NULL DEFAULT false,
    "created_at" bigint NOT NULL,
    "updated_at" bigint NOT NULL,
    "deleted_at" bigint NOT NULL DEFAULT 0
);