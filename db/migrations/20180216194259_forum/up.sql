CREATE TABLE forum_tags (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  color VARCHAR(16) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_forum_tags_name ON forum_tags (name);

CREATE TABLE forum_catalogs (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  summary VARCHAR(800) NOT NULL,
  icon VARCHAR(16) NOT NULL,
  color VARCHAR(16) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_forum_catalogs_title ON forum_catalogs (title);

CREATE TABLE forum_topics (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  type VARCHAR(8) NOT NULL,
  user_id BIGINT NOT NULL REFERENCES users,
  catalog_id BIGINT NOT NULL REFERENCES forum_catalogs,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_forum_topics_title ON forum_topics (title);

CREATE TABLE forum_topics_tags (
  id BIGSERIAL PRIMARY KEY,
  tag_id BIGINT NOT NULL REFERENCES forum_tags,
  topic_id BIGINT NOT NULL REFERENCES forum_topics
);
CREATE INDEX idx_forum_topics_tags_ids ON forum_topics_tags (tag_id, topic_id);

CREATE TABLE forum_posts (
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  type VARCHAR(8) NOT NULL,
  user_id BIGINT NOT NULL REFERENCES users,
  topic_id BIGINT NOT NULL REFERENCES forum_topics,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
