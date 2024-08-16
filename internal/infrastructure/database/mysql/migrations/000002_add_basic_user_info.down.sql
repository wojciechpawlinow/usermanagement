DROP TABLE IF EXISTS addresses;

ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN last_name,
DROP COLUMN phone_number;