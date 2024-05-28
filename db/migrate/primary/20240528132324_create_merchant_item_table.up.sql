CREATE TYPE "product_category" AS ENUM ('Beverage', 'Food', 'Snack', 'Condiments', 'Additions');

CREATE TABLE IF NOT EXISTS "products" (
  id uuid NOT NULL,
  user_id uuid NOT NULL,
  merchant_id uuid NOT NULL,
  name varchar(30) NOT NULL,
  product_category product_category NOT NULL,
  price numeric(10,2) NOT NULL,
  image_url varchar(255) NOT NULL,
  created_at timestamp NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "admin_details" ("user_id") ON DELETE CASCADE,
  FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id") ON DELETE CASCADE
  );
