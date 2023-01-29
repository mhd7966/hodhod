-- +migrate Up
CREATE TABLE project_alert_groups (
  id SERIAL PRIMARY KEY,
  name TEXT,
  user_id INT,
  project_id INT,
  service TEXT,
  alert_group_id INT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE project_alert_groups;