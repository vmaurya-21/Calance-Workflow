-- Drop indexes first
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_github_id;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop users table
DROP TABLE IF EXISTS users;
