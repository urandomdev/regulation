-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "password" bytea NOT NULL,
  "nickname" character varying NOT NULL,
  "custody_account_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_users_custody_account" FOREIGN KEY ("custody_account_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create "plaid_items" table
CREATE TABLE "plaid_items" (
  "id" uuid NOT NULL,
  "plaid_item_id" character varying NOT NULL,
  "access_token" character varying NOT NULL,
  "institution_name" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "plaid_items_users_plaid_items" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "plaid_items_plaid_item_id_key" to table: "plaid_items"
CREATE UNIQUE INDEX "plaid_items_plaid_item_id_key" ON "plaid_items" ("plaid_item_id");
-- Create index "plaiditem_plaid_item_id" to table: "plaid_items"
CREATE UNIQUE INDEX "plaiditem_plaid_item_id" ON "plaid_items" ("plaid_item_id");
-- Create index "plaiditem_user_id_is_active" to table: "plaid_items"
CREATE INDEX "plaiditem_user_id_is_active" ON "plaid_items" ("user_id", "is_active");
-- Create "plaid_accounts" table
CREATE TABLE "plaid_accounts" (
  "id" uuid NOT NULL,
  "plaid_account_id" character varying NOT NULL,
  "account_name" character varying NOT NULL,
  "account_type" character varying NOT NULL DEFAULT 'checking',
  "account_subtype" character varying NULL,
  "mask" character varying NULL,
  "current_balance" bigint NOT NULL DEFAULT 0,
  "available_balance" bigint NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "plaid_item_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "plaid_accounts_plaid_items_accounts" FOREIGN KEY ("plaid_item_id") REFERENCES "plaid_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "plaid_accounts_users_plaid_accounts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "plaid_accounts_plaid_account_id_key" to table: "plaid_accounts"
CREATE UNIQUE INDEX "plaid_accounts_plaid_account_id_key" ON "plaid_accounts" ("plaid_account_id");
-- Create index "plaidaccount_plaid_account_id" to table: "plaid_accounts"
CREATE UNIQUE INDEX "plaidaccount_plaid_account_id" ON "plaid_accounts" ("plaid_account_id");
-- Create index "plaidaccount_plaid_item_id" to table: "plaid_accounts"
CREATE INDEX "plaidaccount_plaid_item_id" ON "plaid_accounts" ("plaid_item_id");
-- Create index "plaidaccount_user_id_is_active" to table: "plaid_accounts"
CREATE INDEX "plaidaccount_user_id_is_active" ON "plaid_accounts" ("user_id", "is_active");
-- Create "plaid_sync_cursors" table
CREATE TABLE "plaid_sync_cursors" (
  "id" uuid NOT NULL,
  "cursor" character varying NOT NULL DEFAULT '',
  "last_sync_at" timestamptz NOT NULL,
  "plaid_item_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "plaid_sync_cursors_plaid_items_sync_cursor" FOREIGN KEY ("plaid_item_id") REFERENCES "plaid_items" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "plaid_sync_cursors_plaid_item_id_key" to table: "plaid_sync_cursors"
CREATE UNIQUE INDEX "plaid_sync_cursors_plaid_item_id_key" ON "plaid_sync_cursors" ("plaid_item_id");
-- Create index "plaidsynccursor_plaid_item_id" to table: "plaid_sync_cursors"
CREATE UNIQUE INDEX "plaidsynccursor_plaid_item_id" ON "plaid_sync_cursors" ("plaid_item_id");
-- Create "transactions" table
CREATE TABLE "transactions" (
  "id" uuid NOT NULL,
  "plaid_transaction_id" character varying NOT NULL,
  "amount" bigint NOT NULL,
  "date" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "merchant_name" character varying NULL,
  "category" character varying NOT NULL DEFAULT 'Misc',
  "plaid_categories" jsonb NULL,
  "pending" boolean NOT NULL DEFAULT false,
  "payment_channel" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "plaid_account_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "transactions_plaid_accounts_transactions" FOREIGN KEY ("plaid_account_id") REFERENCES "plaid_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "transaction_category" to table: "transactions"
CREATE INDEX "transaction_category" ON "transactions" ("category");
-- Create index "transaction_pending" to table: "transactions"
CREATE INDEX "transaction_pending" ON "transactions" ("pending");
-- Create index "transaction_plaid_account_id_date" to table: "transactions"
CREATE INDEX "transaction_plaid_account_id_date" ON "transactions" ("plaid_account_id", "date");
-- Create index "transaction_plaid_transaction_id" to table: "transactions"
CREATE UNIQUE INDEX "transaction_plaid_transaction_id" ON "transactions" ("plaid_transaction_id");
-- Create index "transactions_plaid_transaction_id_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_plaid_transaction_id_key" ON "transactions" ("plaid_transaction_id");
-- Create "virtual_accounts" table
CREATE TABLE "virtual_accounts" (
  "id" uuid NOT NULL,
  "type" character varying NOT NULL,
  "name" character varying NOT NULL,
  "dollars" bigint NOT NULL,
  "cents" bigint NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "virtual_accounts_users_accounts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "virtual_account_transactions" table
CREATE TABLE "virtual_account_transactions" (
  "id" uuid NOT NULL,
  "adjusted_dollars" bigint NOT NULL,
  "adjusted_cents" bigint NOT NULL,
  "virtual_account_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "virtual_account_transactions_virtual_accounts_transactions" FOREIGN KEY ("virtual_account_id") REFERENCES "virtual_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
