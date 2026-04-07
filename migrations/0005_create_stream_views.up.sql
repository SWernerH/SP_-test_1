CREATE TABLE "stream_views" (
  "id" bigserial PRIMARY KEY,
  "stream_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "joined_at" timestamp,
  "left_at" timestamp,
  FOREIGN KEY ("stream_id") REFERENCES "streams" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
