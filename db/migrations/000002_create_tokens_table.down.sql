-- Drop indexes first
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_deleted_at;

-- Drop tokens table
DROP TABLE IF EXISTS tokens;
