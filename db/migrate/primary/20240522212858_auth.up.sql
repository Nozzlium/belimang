CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID NOT NULL,
  "username" VARCHAR(30) NOT NULL,
  PRIMARY KEY ("id"), 
  UNIQUE ("username")
);

CREATE TABLE IF NOT EXISTS "admin_details" (
  "user_id" uuid not null,
  "email" varchar(255) not null,
  "password" varchar(100) not null,
  primary key ("user_id"),
  foreign key ("user_id") references "users" ("id") on delete cascade,
  unique ("email")
);

CREATE TABLE IF NOT EXISTS "user_details" (
  "user_id" uuid not null,
  "email" varchar(255) not null,
  "password" varchar(100) not null,
  PRIMARY KEY ("user_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") on delete cascade,
  unique ("email")
);
