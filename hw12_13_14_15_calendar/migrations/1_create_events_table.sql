-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS "events";
CREATE TABLE "public"."events" (
    "id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "title" varchar NOT NULL,
    "date_start" timestamptz NOT NULL,
    "date_end" timestamptz,
    "user_id" uuid NULL,
    "description" text,
    "date_notification" timestamptz,
    CONSTRAINT "events_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "valid_period"            CHECK (date_start < date_end),
    CONSTRAINT "valid_date_notification" CHECK (date_notification <= date_start)
) WITH (oids = false);

CREATE INDEX events_date_start_idx ON events (date_start);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX events_date_start_idx;

DROP TABLE events;
-- +goose StatementEnd