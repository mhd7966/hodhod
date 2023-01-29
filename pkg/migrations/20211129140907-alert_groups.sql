-- +migrate Up
CREATE TABLE alert_groups (
  id SERIAL PRIMARY KEY,
  name TEXT,
  user_id INT,
  is_default BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE alert_groups;