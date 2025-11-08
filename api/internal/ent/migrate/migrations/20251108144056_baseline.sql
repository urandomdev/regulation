-- Create "push_subscriptions" table
CREATE TABLE "push_subscriptions" (
  "id" uuid NOT NULL,
  "endpoint" character varying NOT NULL,
  "p256dh" character varying NOT NULL,
  "auth" character varying NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "push_subscriptions_users_push_subscriptions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "pushsubscription_endpoint" to table: "push_subscriptions"
CREATE UNIQUE INDEX "pushsubscription_endpoint" ON "push_subscriptions" ("endpoint");
-- Create index "pushsubscription_user_id_active" to table: "push_subscriptions"
CREATE INDEX "pushsubscription_user_id_active" ON "push_subscriptions" ("user_id", "active");
