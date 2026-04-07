CREATE TABLE "donations" (
  "id" bigserial PRIMARY KEY,
  "donor_id" bigint NOT NULL,
  "streamer_id" bigint NOT NULL,
  "stream_id" bigint,
  "amount" decimal(10,2) NOT NULL,
  "created_at" timestamp,
  FOREIGN KEY ("donor_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("streamer_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("stream_id") REFERENCES "streams" ("id")
);
