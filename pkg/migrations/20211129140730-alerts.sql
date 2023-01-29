-- +migrate Up
CREATE TABLE alerts (
  id SERIAL PRIMARY KEY,
  message TEXT,
  severity INT,
  parent_id INT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE alerts;