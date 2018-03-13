CREATE TABLE users (
  id                 BIGSERIAL PRIMARY KEY,
  name               VARCHAR(32)                 NOT NULL,
  email              VARCHAR(255)                NOT NULL,
  uid                VARCHAR(36)                 NOT NULL,
  password           bytea,
  provider_id        VARCHAR(255)                NOT NULL,
  provider_type      VARCHAR(32)                 NOT NULL,
  logo               VARCHAR(255),
  sign_in_count      INT                         NOT NULL DEFAULT 0,
  current_sign_in_at TIMESTAMP WITH TIME ZONE,
  current_sign_in_ip INET,
  last_sign_in_at    TIMESTAMP WITH TIME ZONE,
  last_sign_in_ip    INET,
  confirmed_at       TIMESTAMP WITH TIME ZONE,
  locked_at          TIMESTAMP WITH TIME ZONE,
  created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_users_uid
  ON users (uid);
CREATE UNIQUE INDEX idx_users_email
  ON users (email);
CREATE UNIQUE INDEX idx_users_provider_id_type
  ON users (provider_id, provider_type);
CREATE INDEX idx_users_name
  ON users (name);
CREATE INDEX idx_users_provider_id
  ON users (provider_id);
CREATE INDEX idx_users_provider_type
  ON users (provider_type);


CREATE TABLE logs (
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT                      NOT NULL REFERENCES users,
  ip         INET                        NOT NULL,
  message    VARCHAR(255)                NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
  id            BIGSERIAL PRIMARY KEY,
  name          VARCHAR(32)                 NOT NULL,
  resource_id   BIGINT NOT NULL,
  resource_type VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMP WITH TIME ZONE NOT NULL
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
  exp        DATE                        NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_policies
  ON policies (user_id, role_id);
