CREATE TABLE "streams" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "description" text,
  "is_live" boolean,
  "started_at" timestamp,
  "ended_at" timestamp,
  "created_at" timestamp,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
