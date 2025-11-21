# Database Migrations Guide

This guide explains how to work with database migrations for the Calance Workflow project.

## Overview

We use **golang-migrate** for managing database schema migrations. Migrations are stored as SQL files in the `db/migrations/` directory.

## Installation

### Install golang-migrate CLI

**Windows (PowerShell):**
```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Linux/macOS:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify installation:
```bash
migrate -version
```

## Migration File Structure

Migration files follow this naming convention:
```
{version}_{description}.{up|down}.sql
```

Examples:
- `000001_create_users_table.up.sql` - Apply migration
- `000001_create_users_table.down.sql` - Rollback migration

**Current Migrations:**
1. `000001_create_users_table` - Creates users table with GitHub OAuth fields
2. `000002_create_tokens_table` - Creates tokens table with one-to-one relationship to users

## Running Migrations

### Setup Database Connection

Set your database URL as an environment variable or use it directly:
```bash
# Windows PowerShell
$env:DATABASE_URL="postgresql://postgres:postgres@localhost:5432/calance_workflow?sslmode=disable"

# Linux/macOS
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/calance_workflow?sslmode=disable"
```

### Apply All Migrations (Up)

```bash
migrate -path db/migrations -database $env:DATABASE_URL up
```

Or with inline database URL:
```bash
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/calance_workflow?sslmode=disable" up
```

### Apply Specific Number of Migrations

```bash
# Apply next migration
migrate -path db/migrations -database $env:DATABASE_URL up 1

# Apply next 2 migrations
migrate -path db/migrations -database $env:DATABASE_URL up 2
```

### Rollback Migrations (Down)

```bash
# Rollback last migration
migrate -path db/migrations -database $env:DATABASE_URL down 1

# Rollback all migrations
migrate -path db/migrations -database $env:DATABASE_URL down
```

### Check Migration Status

```bash
migrate -path db/migrations -database $env:DATABASE_URL version
```

### Force Migration Version

If a migration fails and leaves the database in a "dirty" state:
```bash
# Force to a specific version (use with caution!)
migrate -path db/migrations -database $env:DATABASE_URL force <version>

# Example: Force to version 1
migrate -path db/migrations -database $env:DATABASE_URL force 1
```

## Creating New Migrations

### Using migrate CLI

```bash
# Create a new migration
migrate create -ext sql -dir db/migrations -seq create_new_table
```

This creates two files:
- `000003_create_new_table.up.sql`
- `000003_create_new_table.down.sql`

### Manual Creation

1. Determine the next version number (e.g., `000003`)
2. Create two files in `db/migrations/`:
   - `000003_description.up.sql` - Changes to apply
   - `000003_description.down.sql` - How to revert changes

### Migration Best Practices

1. **Always create both up and down migrations** - Ensure you can rollback
2. **Keep migrations atomic** - One logical change per migration
3. **Use idempotent operations** - Use `IF NOT EXISTS`, `IF EXISTS` clauses
4. **Never modify existing migrations** - Once applied to any environment, create a new migration instead
5. **Test thoroughly** - Test both up and down migrations in development
6. **Use transactions** - PostgreSQL wraps each migration in a transaction by default

## Example Migration Files

### Up Migration Example
```sql
-- 000003_add_user_role.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user';
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

### Down Migration Example
```sql
-- 000003_add_user_role.down.sql
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
```

## Database Schema Verification

### Connect to PostgreSQL

Using Docker:
```bash
docker exec -it calance-workflow-db psql -U postgres -d calance_workflow
```

### Useful PostgreSQL Commands

```sql
-- List all tables
\dt

-- Describe a specific table
\d users
\d tokens

-- View table constraints
SELECT constraint_name, constraint_type 
FROM information_schema.table_constraints 
WHERE table_name = 'tokens';

-- Check foreign keys
SELECT
    tc.table_name, 
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name='tokens';
```

## Development Workflow

### Option 1: Using golang-migrate (Recommended for Production)

1. Start Docker Compose:
   ```bash
   docker-compose up -d
   ```

2. Run migrations:
   ```bash
   migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/calance_workflow?sslmode=disable" up
   ```

3. Start your application:
   ```bash
   go run cmd/server/main.go
   ```

### Option 2: Using GORM AutoMigrate (Development Only)

The application uses GORM's AutoMigrate in `internal/database/database.go` for development convenience. This will automatically create/update tables when you run the application:

```bash
docker-compose up -d
go run cmd/server/main.go
```

> **Note:** For production, always use golang-migrate with SQL migration files for better control and explicit schema versioning.

## Troubleshooting

### "dirty database version" Error

This occurs when a migration fails partway through:

1. Check the current version:
   ```bash
   migrate -path db/migrations -database $env:DATABASE_URL version
   ```

2. Manually fix the database or rollback to a clean state

3. Force the version to match reality:
   ```bash
   migrate -path db/migrations -database $env:DATABASE_URL force <version>
   ```

### Connection Refused

Ensure PostgreSQL is running:
```bash
docker-compose ps
```

If not running:
```bash
docker-compose up -d
```

### Migration Files Not Found

Ensure you're running the command from the project root directory where `db/migrations/` exists.

## CI/CD Integration

For automated deployments, add migration step to your CI/CD pipeline:

```yaml
# Example GitHub Actions workflow
- name: Run Database Migrations
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}
  run: |
    migrate -path db/migrations -database "$DATABASE_URL" up
```

## Additional Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [GORM Documentation](https://gorm.io/docs/)
