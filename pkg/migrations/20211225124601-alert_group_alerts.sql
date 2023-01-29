
-- +migrate Up
CREATE TABLE alert_group_alerts (
  id SERIAL PRIMARY KEY,
  alert_group_id INT,
  alert_id INT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE alert_group_alerts;