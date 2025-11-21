-- Enable UUID extension for PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    github_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    avatar_url TEXT,
    name VARCHAR(255),
    bio TEXT,
    location VARCHAR(255),
    company VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;

-- Add comment
COMMENT ON TABLE users IS 'Stores GitHub user information from OAuth authentication';
