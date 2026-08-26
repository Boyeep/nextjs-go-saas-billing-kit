DROP TRIGGER IF EXISTS users_create_personal_workspace ON users;
DROP FUNCTION IF EXISTS create_personal_workspace();
DROP TABLE IF EXISTS usage_counters;
DROP TABLE IF EXISTS billing_events;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS workspace_memberships;
DROP TABLE IF EXISTS workspaces;
DROP TYPE IF EXISTS membership_role;

