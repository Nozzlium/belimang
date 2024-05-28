CREATE TYPE "merchant_category" AS ENUM ('SmallRestaurant', 'MediumRestaurant', 'LargeRestaurant', 'MerchandiseRestaurant', 'BoothKiosk', 'ConvenienceStore');

CREATE TABLE IF NOT EXISTS "merchants" (
  id uuid NOT NULL,
  user_id uuid NOT NULL,
  name varchar(30) NOT NULL,
  merchant_category MERCHANT_CATEGORY NOT NULL,
  image_url varchar(255) NOT NULL,
  latitude float NOT NULL,
  longitude float NOT NULL,
  created_at timestamp NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "admin_details" ("user_id") ON DELETE CASCADE
);
