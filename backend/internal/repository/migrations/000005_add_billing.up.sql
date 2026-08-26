DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'membership_role') THEN
    CREATE TYPE membership_role AS ENUM ('owner', 'admin', 'member');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS workspaces (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(120) NOT NULL,
  slug VARCHAR(80) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workspace_memberships (
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role membership_role NOT NULL DEFAULT 'member',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (workspace_id, user_id)
);

CREATE TABLE IF NOT EXISTS subscriptions (
  workspace_id UUID PRIMARY KEY REFERENCES workspaces(id) ON DELETE CASCADE,
  provider VARCHAR(20) NOT NULL DEFAULT 'fake',
  provider_customer_id TEXT,
  provider_subscription_id TEXT UNIQUE,
  plan VARCHAR(30) NOT NULL DEFAULT 'free',
  status VARCHAR(30) NOT NULL DEFAULT 'trialing',
  current_period_end TIMESTAMPTZ,
  grace_period_end TIMESTAMPTZ,
  cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing_events (
  provider_event_id TEXT PRIMARY KEY,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS usage_counters (
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  metric VARCHAR(80) NOT NULL,
  period_start DATE NOT NULL,
  quantity BIGINT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  PRIMARY KEY (workspace_id, metric, period_start)
);

CREATE OR REPLACE FUNCTION create_personal_workspace() RETURNS TRIGGER AS $$
DECLARE workspace_uuid UUID;
BEGIN
  INSERT INTO workspaces (name, slug)
  VALUES (COALESCE(NULLIF(NEW.full_name, ''), NEW.username) || '''s workspace', 'ws-' || NEW.id)
  RETURNING id INTO workspace_uuid;
  INSERT INTO workspace_memberships (workspace_id, user_id, role)
  VALUES (workspace_uuid, NEW.id, 'owner');
  INSERT INTO subscriptions (workspace_id) VALUES (workspace_uuid);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS users_create_personal_workspace ON users;
CREATE TRIGGER users_create_personal_workspace
AFTER INSERT ON users FOR EACH ROW EXECUTE FUNCTION create_personal_workspace();

INSERT INTO workspaces (name, slug)
SELECT COALESCE(NULLIF(full_name, ''), username) || '''s workspace', 'ws-' || id FROM users
ON CONFLICT (slug) DO NOTHING;

INSERT INTO workspace_memberships (workspace_id, user_id, role)
SELECT w.id, u.id, 'owner' FROM users u JOIN workspaces w ON w.slug = 'ws-' || u.id
ON CONFLICT DO NOTHING;

INSERT INTO subscriptions (workspace_id)
SELECT id FROM workspaces ON CONFLICT DO NOTHING;

