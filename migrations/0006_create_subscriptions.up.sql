CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "subscriber_id" bigint NOT NULL,
  "streamer_id" bigint NOT NULL,
  "amount" decimal(10,2) NOT NULL,
  "created_at" timestamp,
  FOREIGN KEY ("subscriber_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("streamer_id") REFERENCES "users" ("id")
);
