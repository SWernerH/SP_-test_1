CREATE TABLE "follows" (
  "id" bigserial PRIMARY KEY,
  "follower_id" bigint NOT NULL,
  "streamer_id" bigint NOT NULL,
  "created_at" timestamp,
  FOREIGN KEY ("follower_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("streamer_id") REFERENCES "users" ("id")
);
