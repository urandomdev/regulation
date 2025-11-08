-- Create "items" table
CREATE TABLE "items" (
  "id" uuid NOT NULL,
  "plaid_id" character varying NOT NULL,
  "access_token" character varying NOT NULL,
  "institution_name" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "items_users_items" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "item_plaid_id" to table: "items"
CREATE UNIQUE INDEX "item_plaid_id" ON "items" ("plaid_id");
-- Create index "item_user_id_is_active" to table: "items"
CREATE INDEX "item_user_id_is_active" ON "items" ("user_id", "is_active");
-- Create index "items_plaid_id_key" to table: "items"
CREATE UNIQUE INDEX "items_plaid_id_key" ON "items" ("plaid_id");
-- Create "accounts" table
CREATE TABLE "accounts" (
  "id" uuid NOT NULL,
  "plaid_id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "type" character varying NOT NULL DEFAULT 'checking',
  "subtype" character varying NULL,
  "mask" character varying NULL,
  "current_balance" bigint NOT NULL DEFAULT 0,
  "available_balance" bigint NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "item_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "accounts_items_accounts" FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "accounts_users_accounts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "account_item_id" to table: "accounts"
CREATE INDEX "account_item_id" ON "accounts" ("item_id");
-- Create index "account_plaid_id" to table: "accounts"
CREATE UNIQUE INDEX "account_plaid_id" ON "accounts" ("plaid_id");
-- Create index "account_user_id_is_active" to table: "accounts"
CREATE INDEX "account_user_id_is_active" ON "accounts" ("user_id", "is_active");
-- Create index "accounts_plaid_id_key" to table: "accounts"
CREATE UNIQUE INDEX "accounts_plaid_id_key" ON "accounts" ("plaid_id");
-- Modify "transactions" table
ALTER TABLE "transactions" DROP COLUMN "plaid_transaction_id", DROP COLUMN "plaid_account_id", ADD COLUMN "plaid_id" character varying NOT NULL, ADD COLUMN "account_id" uuid NOT NULL, ADD CONSTRAINT "transactions_accounts_transactions" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "transaction_account_id_date" to table: "transactions"
CREATE INDEX "transaction_account_id_date" ON "transactions" ("account_id", "date");
-- Create index "transaction_plaid_id" to table: "transactions"
CREATE UNIQUE INDEX "transaction_plaid_id" ON "transactions" ("plaid_id");
-- Create index "transactions_plaid_id_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_plaid_id_key" ON "transactions" ("plaid_id");
-- Create "sync_cursors" table
CREATE TABLE "sync_cursors" (
  "id" uuid NOT NULL,
  "cursor" character varying NOT NULL DEFAULT '',
  "last_sync_at" timestamptz NOT NULL,
  "item_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "sync_cursors_items_sync_cursor" FOREIGN KEY ("item_id") REFERENCES "items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "sync_cursors_item_id_key" to table: "sync_cursors"
CREATE UNIQUE INDEX "sync_cursors_item_id_key" ON "sync_cursors" ("item_id");
-- Create index "synccursor_item_id" to table: "sync_cursors"
CREATE UNIQUE INDEX "synccursor_item_id" ON "sync_cursors" ("item_id");
-- Drop "plaid_accounts" table
DROP TABLE "plaid_accounts";
-- Drop "plaid_sync_cursors" table
DROP TABLE "plaid_sync_cursors";
-- Drop "plaid_items" table
DROP TABLE "plaid_items";
-- Drop "virtual_account_transactions" table
DROP TABLE "virtual_account_transactions";
-- Drop "virtual_accounts" table
DROP TABLE "virtual_accounts";
