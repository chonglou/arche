CREATE TABLE donate_projects (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  type VARCHAR(8) NOT NULL,
  methods TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_donate_projects_title ON donate_projects (title);
CREATE INDEX idx_donate_projects_type ON donate_projects (type);
