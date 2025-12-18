# Database Migration Guide

## Overview

This project uses **[golang-migrate](https://github.com/golang-migrate/migrate)** for database schema versioning and migrations. This provides:

- ✅ **Versioned migrations** - Track migration history in database
- ✅ **Up/Down migrations** - Forward migration and rollback capability
- ✅ **SQL-based** - Write migrations in pure SQL
- ✅ **CLI tool** - Manage migrations from command line or code
- ✅ **Production-ready** - Industry standard tool

## Migration Files

Migrations are stored in the `migrations/` directory with the naming convention:

```
{version}_{name}.up.sql    - Forward migration
{version}_{name}.down.sql  - Rollback migration
```

Example:
```
migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_add_user_avatar.up.sql
└── 000002_add_user_avatar.down.sql
```

## Quick Start

### Automatic Migrations (Development)

Migrations run automatically when the application starts:

```bash
make run
# or
docker-compose up
```

The application will:
1. Connect to the database
2. Check for pending migrations
3. Apply all pending migrations
4. Start the server

### Manual Migration Commands

#### Apply All Pending Migrations
```bash
make migrate-up
```

#### Rollback Last Migration
```bash
make migrate-down
```

#### Check Migration Status
```bash
make migrate-status
```

#### Show Current Version
```bash
make migrate-version
```

#### Create New Migration
```bash
make migrate-create
# Enter migration name when prompted, e.g., "add_user_avatar"
```

This creates two files:
- `migrations/000002_add_user_avatar.up.sql`
- `migrations/000002_add_user_avatar.down.sql`

## Migration Workflow

### 1. Creating a New Migration

```bash
# Create migration files
make migrate-create
# Enter: add_user_avatar

# This creates:
# migrations/000002_add_user_avatar.up.sql
# migrations/000002_add_user_avatar.down.sql
```

### 2. Write Up Migration

Edit `migrations/000002_add_user_avatar.up.sql`:

```sql
-- Add avatar_url column to users table
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(500);

-- Create index for faster queries
CREATE INDEX idx_users_avatar ON users(avatar_url) WHERE avatar_url IS NOT NULL;
```

### 3. Write Down Migration

Edit `migrations/000002_add_user_avatar.down.sql`:

```sql
-- Remove index
DROP INDEX IF EXISTS idx_users_avatar;

-- Remove avatar_url column
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
```

### 4. Apply Migration

```bash
# Apply the new migration
make migrate-up

# Check status
make migrate-status
```

### 5. Rollback if Needed

```bash
# Rollback last migration
make migrate-down

# Check status
make migrate-status
```

## Migration Best Practices

### ✅ Do's

1. **Always write down migrations**
   - Every up migration must have a corresponding down migration
   - Test rollback before deploying to production

2. **Use transactions for data migrations**
   ```sql
   BEGIN;
   
   -- Your migration here
   UPDATE users SET status = 'active' WHERE status IS NULL;
   
   COMMIT;
   ```

3. **Make migrations idempotent**
   ```sql
   -- Good: Can run multiple times safely
   CREATE TABLE IF NOT EXISTS new_table (...);
   ALTER TABLE users ADD COLUMN IF NOT EXISTS new_col VARCHAR(100);
   
   -- Bad: Fails if run twice
   CREATE TABLE new_table (...);
   ALTER TABLE users ADD COLUMN new_col VARCHAR(100);
   ```

4. **Add indexes for foreign keys**
   ```sql
   ALTER TABLE courses ADD COLUMN category_id UUID;
   CREATE INDEX idx_courses_category ON courses(category_id);
   ```

5. **Use constraints for data integrity**
   ```sql
   ALTER TABLE users ADD COLUMN age INTEGER CHECK (age >= 0 AND age <= 150);
   ```

### ❌ Don'ts

1. **Never modify existing migrations**
   - Once applied to any environment, migrations are immutable
   - Create a new migration to fix issues

2. **Avoid destructive operations without backups**
   - Dropping tables/columns in production is dangerous
   - Consider soft deletes or archiving first

3. **Don't mix DDL and DML**
   - Keep schema changes (DDL) and data changes (DML) separate
   - Makes rollback easier

4. **Don't use database-specific features without consideration**
   - Stay portable if you might switch databases
   - Document database-specific features

## Migration States

### Clean State
```
Current Version: 1
Status: clean
```
- All migrations applied successfully
- Database is in sync with migration files

### Dirty State
```
Current Version: 1
Status: dirty (migration failed, needs manual intervention)
```
- A migration failed partway through
- Database may be in inconsistent state
- Requires manual intervention

#### Recovering from Dirty State

1. **Check the database state**
   ```bash
   make db-shell
   # Inspect tables and data
   ```

2. **Fix the issue**
   - If migration partially applied: Complete it manually or undo it
   - If migration not applied: Fix the migration SQL

3. **Force version** (use with caution)
   ```bash
   # If you manually fixed the database
   make migrate-force
   # Enter the correct version number
   ```

4. **Continue migrations**
   ```bash
   make migrate-up
   ```

## Advanced Usage

### Using the Migration CLI Directly

```bash
# Set database URL
export DATABASE_URL="postgresql://lms_user:lms_password@localhost:5432/lms_db?sslmode=disable"

# Apply migrations
go run src/cmd/migrate/main.go up

# Rollback
go run src/cmd/migrate/main.go down

# Check status
go run src/cmd/migrate/main.go status

# Force version (emergency only)
go run src/cmd/migrate/main.go force 1
```

### Using migrate CLI Tool

If you prefer the standalone CLI:

```bash
# Install (already done via homebrew)
brew install golang-migrate

# Apply migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback
migrate -path migrations -database "$DATABASE_URL" down

# Create migration
migrate create -ext sql -dir migrations -seq migration_name
```

## Production Deployment

### Option 1: Automatic on Startup (Current)

The application runs migrations automatically on startup:

```go
// In main.go
if err := runMigrations(&cfg.Database); err != nil {
    log.Printf("⚠️  Migration warning: %v", err)
}
```

**Pros:**
- Simple, no extra steps
- Works for small teams and deployments

**Cons:**
- Application waits for migrations
- Not ideal for zero-downtime deployments

### Option 2: Separate Migration Job

For production, run migrations as a separate step:

```bash
# 1. Run migrations first
go run src/cmd/migrate/main.go up

# 2. Then start application
go run src/cmd/api/main.go
```

**Kubernetes Example:**
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: lms-migrate
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: lms-api:latest
        command: ["go", "run", "src/cmd/migrate/main.go", "up"]
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
      restartPolicy: OnFailure
```

### Option 3: CI/CD Pipeline

Run migrations in your deployment pipeline:

```yaml
# GitHub Actions example
- name: Run Database Migrations
  run: |
    go run src/cmd/migrate/main.go up
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}

- name: Deploy Application
  run: |
    kubectl apply -f k8s/
```

## Migration History

### Current Schema (Version 1)

**Migration:** `000001_initial_schema`

**Tables:**
- `users` - User accounts with roles
- `courses` - Course catalog
- `modules` - Course modules
- `lessons` - Lesson records
- `lesson_versions` - Lesson content with versioning
- `course_enrollments` - Student enrollments
- `audit_logs` - Audit trail

**Features:**
- UUID primary keys
- Soft deletes (deleted_at)
- Foreign key constraints
- Performance indexes
- Check constraints for enums
- JSONB for audit logs

## Troubleshooting

### Migration Fails with "dirty database"

```bash
# 1. Check current state
make migrate-status

# 2. Inspect database
make db-shell

# 3. Fix manually, then force version
make migrate-force
```

### Can't Find Migration Files

Ensure you're running commands from project root:

```bash
cd /path/to/go-learning-management-system
make migrate-up
```

### Database Connection Issues

Check your `.env` file:

```bash
DATABASE_URL=postgresql://lms_user:lms_password@localhost:5432/lms_db?sslmode=disable
```

### Migration Already Applied

```bash
$ make migrate-up
no change
```

This is normal - means all migrations are already applied.

## Schema Version Table

Migrations are tracked in the `schema_migrations` table:

```sql
-- View migration history
SELECT * FROM schema_migrations;
```

Output:
```
 version | dirty 
---------+-------
       1 | f
```

## Testing Migrations

### Test Up Migration

```bash
# Start fresh database
docker-compose down -v
docker-compose up -d postgres

# Apply migration
make migrate-up

# Verify schema
make db-shell
\dt
\d users
```

### Test Down Migration

```bash
# Rollback
make migrate-down

# Verify tables removed
make db-shell
\dt
```

### Test Idempotency

```bash
# Apply twice - should succeed both times
make migrate-up
make migrate-up
```

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)

## Support

For migration issues:
1. Check migration status: `make migrate-status`
2. Review migration files in `migrations/`
3. Check application logs
4. Consult this guide

---

**Last Updated:** December 17, 2025  
**Schema Version:** 1  
**Migration Tool:** golang-migrate v4.19.1
