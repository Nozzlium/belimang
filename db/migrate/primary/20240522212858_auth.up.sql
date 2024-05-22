CREATE TABLE IF NOT EXISTS "admins" (
  "id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(100) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("email")
)

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(100) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("email")
)

CREATE TABLE IF NOT EXISTS "usernames" (
  "user_id" uuid NOT NULL,
  "username" varchar(30) NOT NULL,
  PRIMARY KEY ("user_id", "username"),
  FOREIGN KEY ("id") REFERENCES "users" ("id") ON DELETE CASCADE,
  UNIQUE ("username")
)
