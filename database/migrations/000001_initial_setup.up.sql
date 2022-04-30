create type enum_approval_status as enum ('pending','approved','rejected');
create type enum_payment_status as enum ('paid','unpaid');

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" text NOT NULL,
  "mobile" text,
  "address" text,
  "email" text UNIQUE NOT NULL,
  "hashed_password" text NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "is_active" bool NOT NULL DEFAULT 'true',
  "created_by" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated_by" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "ip_from" text NOT NULL,
  "user_agent" text NOT NULL
);

CREATE TABLE "borrowers" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "loan_id" bigint NOT NULL,
  "is_active" bool NOT NULL DEFAULT 'true',
  "created_by" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated_by" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "ip_from" text NOT NULL,
  "user_agent" text NOT NULL
);

CREATE TABLE "loans" (
  "id" BIGSERIAL PRIMARY KEY,
  "amount" numeric NOT NULL,
  "amount_need_to_pay" numeric NOT NULL,
  "term" integer NOT NULL,
  "approval_status" enum_approval_status NOT NULL,
  "is_active" bool NOT NULL DEFAULT 'true',
  "repayment_status" enum_payment_status NOT NULL,
  "created_by" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated_by" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "ip_from" text NOT NULL,
  "user_agent" text NOT NULL
);

CREATE TABLE "payments" (
  "id" BIGSERIAL PRIMARY KEY,
  "loan_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  "created_by" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated_by" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "ip_from" text NOT NULL,
  "user_agent" text NOT NULL
);

CREATE INDEX ON "users" ("mobile");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "borrowers" ("user_id");

CREATE INDEX ON "borrowers" ("loan_id");

CREATE INDEX ON "borrowers" ("user_id", "loan_id");

CREATE INDEX ON "payments" ("user_id");

CREATE INDEX ON "payments" ("loan_id");

COMMENT ON COLUMN "payments"."loan_id" IS 'payment made against loan';

COMMENT ON COLUMN "payments"."user_id" IS 'payment made by';

ALTER TABLE "borrowers" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "borrowers" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");