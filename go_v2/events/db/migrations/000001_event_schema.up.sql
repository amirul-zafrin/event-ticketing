CREATE TABLE "events" (
  "id" integer PRIMARY KEY,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "name" varchar,
  "date" timestamp,
  "location" varchar,
  "capacity" integer,
  "seats" JSON
);

CREATE TABLE "prices" (
  "id" integer PRIMARY KEY,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "class" varchar,
  "price" float,
  "event" int
);

ALTER TABLE "prices" ADD FOREIGN KEY ("event") REFERENCES "events" ("id");
