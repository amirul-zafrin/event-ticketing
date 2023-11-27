CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamp NOT NULL DEFAULT 'now()',
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "name" varchar NOT NULL,
  "date" timestamp,
  "location" varchar NOT NULL,
  "capacity" integer NOT NULL,
  "seats" JSONB
);

CREATE TABLE "prices" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamp NOT NULL DEFAULT 'now()',
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "class" varchar NOT NULL,
  "price" float NOT NULL,
  "event" int NOT NULL
);

ALTER TABLE "prices" ADD FOREIGN KEY ("event") REFERENCES "events" ("id");
