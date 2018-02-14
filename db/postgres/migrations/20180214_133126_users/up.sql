CREATE TABLE domains (
  id BIGSERIAL PRIMARY KEY,
  name varchar(64) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_domains_name ON domains(name);

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  domain_id BIGINT NOT NULL REFERENCES users,
  name varchar(64) NOT NULL,
  logo varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  email varchar(128) NOT NULL,
  uid VARCHAR(36) NOT NULL,
  sign_in_count      INT NOT NULL DEFAULT 0,
  current_sign_in_at TIMESTAMP WITHOUT TIME ZONE,
  current_sign_in_ip INET,
  last_sign_in_at    TIMESTAMP WITHOUT TIME ZONE,
  last_sign_in_ip    INET,
  confirmed_at TIMESTAMP WITHOUT TIME ZONE,
  locked_at TIMESTAMP WITHOUT TIME ZONE,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_uid ON users(uid);
CREATE INDEX idx_users_name ON users(name);

CREATE TABLE aliases (
  id BIGSERIAL PRIMARY KEY,
  domain_id BIGINT NOT NULL REFERENCES domains,
  source varchar(128) NOT NULL,
  destination varchar(128) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_aliases_source ON aliases(source);
CREATE INDEX idx_aliases_destination ON aliases(destination);


CREATE TABLE logs (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      NOT NULL REFERENCES users,
  ip         INET                        NOT NULL,
  message    VARCHAR(255)                NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

CREATE TABLE roles (
  id            BIGSERIAL PRIMARY KEY,
  name          VARCHAR(32)                 NOT NULL,
  resource_id   BIGINT,
  resource_type VARCHAR(255),
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_roles_name_resource_type_id
  ON roles (name, resource_type, resource_id);
CREATE INDEX idx_roles_name
  ON roles (name);
CREATE INDEX idx_roles_resource_type
  ON roles (resource_type);

CREATE TABLE policies (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      NOT NULL REFERENCES users,
  role_id    BIGINT                      NOT NULL REFERENCES roles,
  nbf        DATE                        NOT NULL DEFAULT current_date,
  exp        DATE                        NOT NULL DEFAULT '2018-02-14',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_policies
  ON policies (user_id, role_id);
