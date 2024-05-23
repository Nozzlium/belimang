CREATE TABLE IF NOT EXISTS "admins" (
  "id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(100) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("email")
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(100) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("email")
);

CREATE TABLE IF NOT EXISTS "admin_usernames" (
  "admin_id" uuid NOT NULL,
  "username" varchar(30) NOT NULL,
  PRIMARY KEY ("admin_id",  "username"),
  FOREIGN KEY ("admin_id") REFERENCES "admins" ("id") ON DELETE CASCADE,
  UNIQUE ("username")
);

CREATE TABLE IF NOT EXISTS "user_usernames" (
  "user_id" UUID NOT NULL,
  "username" VARCHAR(30) NOT NULL,
  PRIMARY KEY ("user_id", "username"), 
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  UNIQUE ("username")
);
