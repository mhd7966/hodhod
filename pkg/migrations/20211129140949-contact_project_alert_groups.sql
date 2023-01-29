-- +migrate Up
CREATE TABLE contact_project_alert_groups (
  id SERIAL PRIMARY KEY,
  project_alert_group_id INT,
  contact_id INT,
  severity INT,
  call BOOLEAN,
  sms BOOLEAN,
  email BOOLEAN,
  webhook BOOLEAN,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
-- +migrate Down
DROP TABLE contact_project_alert_groups;