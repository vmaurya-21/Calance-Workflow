-- Create tokens table with one-to-one relationship to users
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE, -- UNIQUE constraint ensures one-to-one relationship
    access_token TEXT NOT NULL,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    scope TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraint
    CONSTRAINT fk_tokens_user FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON UPDATE CASCADE 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tokens_deleted_at ON tokens(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);

-- Add comment
COMMENT ON TABLE tokens IS 'Stores OAuth access tokens with one-to-one mapping to users';
COMMENT ON COLUMN tokens.user_id IS 'Unique constraint ensures each user can only have one active token';
