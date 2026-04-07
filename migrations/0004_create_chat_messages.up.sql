CREATE TABLE "chat_messages" (
  "id" bigserial PRIMARY KEY,
  "stream_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "message" text NOT NULL,
  "created_at" timestamp,
  FOREIGN KEY ("stream_id") REFERENCES "streams" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
