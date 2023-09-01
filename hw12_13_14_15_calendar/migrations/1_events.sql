-- Adminer 4.8.1 PostgreSQL 15.4 (Debian 15.4-1.pgdg120+1) dump

DROP TABLE IF EXISTS "events";
CREATE TABLE "public"."events" (
    "id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "title" character(50) NOT NULL,
    "date_start" timestamp NOT NULL,
    "date_end" timestamp,
    "user_id" uuid NULL,
    "description" text NOT NULL,
    "date_notification" timestamp,
    CONSTRAINT "events_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "events" ("id", "title", "date_start", "date_end", "user_id", "description", "date_notification") VALUES
('db223c0d-8623-4ce0-926f-a3c0136cf704',	'test title',	'2023-09-02 09:00:00',	'2023-09-02 11:00:00', 'db223c0d-8623-4ce0-926f-0000136cf704','description',	'2023-09-01 21:00:00');

-- 2023-08-31 19:14:44.025681+00