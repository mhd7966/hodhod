-- +migrate Up
CREATE TABLE logs (
  id SERIAL PRIMARY KEY,
  user_id INT,
  contact_id INT,
  alert_id INT,
  alert_group_id INT,
  project_alert_group_id INT,
  contact_project_alert_group_id INT,
  message TEXT,
  service TEXT,
  project_id INT,
  channel INT,
  status INT,
  driver INT,
  contact_info TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE logs;