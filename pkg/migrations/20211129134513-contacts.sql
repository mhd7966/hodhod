-- +migrate Up
CREATE TABLE contacts (
  id SERIAL PRIMARY KEY,
  user_id INT,
  name TEXT,
  phone_number TEXT,
  email TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE contacts;