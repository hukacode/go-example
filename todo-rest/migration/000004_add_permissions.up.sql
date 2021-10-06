CREATE TABLE IF NOT EXISTS permission (
  id bigserial PRIMARY KEY,
  code text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_permission (
  user_id bigint NOT NULL REFERENCES app_user ON DELETE CASCADE,
  permission_id bigint NOT NULL REFERENCES permission ON DELETE CASCADE,
  PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permission (code) VALUES('task:read'), ('task:write');