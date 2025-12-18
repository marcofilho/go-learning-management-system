# Database Migration System ✅

## Current System: golang-migrate

The project now uses **golang-migrate** for versioned database migrations instead of GORM's AutoMigrate. See [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for complete documentation.

## Previous: GORM Migration

This document describes the historical migration from raw SQL queries (pgx) to GORM ORM.

## Changes Made

### 1. Dependencies Updated
- Added `gorm.io/gorm v1.25.5`
- Added `gorm.io/driver/postgres v1.5.4`
- GORM uses pgx internally as the PostgreSQL driver

### 2. Entity Updates
All domain entities updated with GORM tags:

**User Entity** ([src/internal/domain/entity/user.go](src/internal/domain/entity/user.go))
- Added GORM struct tags (primaryKey, uniqueIndex, type, default)
- Added soft delete support with `gorm.DeletedAt`
- Added `TableName()` method returning "users"
- Auto-generated UUID and timestamps

**Course Entity** ([src/internal/domain/entity/course.go](src/internal/domain/entity/course.go))
- Added GORM struct tags with foreign key to User
- Added Instructor relationship with cascade delete
- Added soft delete support
- Added `TableName()` method returning "courses"

**Enrollment Entity** ([src/internal/domain/entity/enrollment.go](src/internal/domain/entity/enrollment.go))
- Added GORM struct tags with foreign keys to User and Course
- Added User and Course relationships with cascade delete
- Added soft delete support
- Added `TableName()` method returning "enrollments"

### 3. Database Layer
**Connection** ([src/internal/infrastructure/database/postgres.go](src/internal/infrastructure/database/postgres.go))
- Changed from `pgxpool.Pool` to `gorm.DB`
- Updated interface: `GetPool()` → `GetDB()`
- Configured connection pooling (5 idle, 25 max, 1 hour lifetime)
- Added GORM logger configuration

**Migrations** ([src/internal/infrastructure/database/migrations.go](src/internal/infrastructure/database/migrations.go))
- Created `RunMigrations()` function using GORM AutoMigrate
- Automatic schema creation for all entities
- Custom indexes for performance:
  - users: role, is_active
  - courses: instructor_id, status, is_published
  - enrollments: user_id, course_id, status
- Unique constraint on enrollments(user_id, course_id)

### 4. Repository Layer
All three repositories converted from raw SQL to GORM:

**User Repository** ([src/internal/infrastructure/repository/user_repository_postgres.go](src/internal/infrastructure/repository/user_repository_postgres.go))
- `Create()` → `db.Create(&user)`
- `GetByID()` → `db.First(&user, "id = ?", id)`
- `GetByEmail()` → `db.Where("email = ?", email).First(&user)`
- `Update()` → `db.Save(&user)`
- `Delete()` → `db.Delete(&user)` (soft delete)
- `List()` → `db.Order().Limit().Offset().Find(&users)`
- `GetByRole()` → `db.Where("role = ?", role).Find(&users)`

**Course Repository** ([src/internal/infrastructure/repository/course_repository_postgres.go](src/internal/infrastructure/repository/course_repository_postgres.go))
- Similar GORM patterns for CRUD operations
- `GetPublished()` → `db.Where("is_published = ? AND status = ?").Find(&courses)`
- `GetByInstructor()` → `db.Where("instructor_id = ?").Find(&courses)`

**Enrollment Repository** ([src/internal/infrastructure/repository/enrollment_repository_postgres.go](src/internal/infrastructure/repository/enrollment_repository_postgres.go))
- GORM handles nullable fields naturally (CompletedAt, LastAccessAt)
- `GetByUserAndCourse()` → `db.Where("user_id = ? AND course_id = ?").First()`
- `GetActiveByUser()` → `db.Where("user_id = ? AND status = ?").Find()`

### 5. Main Application
**Entry Point** ([src/cmd/api/main.go](src/cmd/api/main.go))
- Changed `pool := db.GetPool()` to `gormDB := db.GetDB()`
- Added migration execution: `database.RunMigrations(gormDB)`
- Updated repository constructors to accept `*gorm.DB`

## Benefits of GORM

### 1. **Productivity**
- Less boilerplate code
- No manual SQL writing
- Auto-handling of NULL values
- Built-in pagination and ordering

### 2. **Type Safety**
- Compile-time checking
- Less prone to SQL injection
- Easier refactoring

### 3. **Features**
- Soft deletes out of the box
- Automatic migrations
- Relationship management
- Hooks and callbacks
- Transaction support

### 4. **Maintainability**
- Cleaner, more readable code
- Self-documenting through struct tags
- Consistent error handling
- Better testing support

## Testing the Application

### 1. Start PostgreSQL
```bash
docker run --name postgres-lms \
  -e POSTGRES_USER=lms_user \
  -e POSTGRES_PASSWORD=lms_password \
  -e POSTGRES_DB=lms_db \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### 2. Run the Application
```bash
go run src/cmd/api/main.go
```

The migrations will run automatically on startup!

### 3. Build for Production
```bash
go build -o bin/api ./src/cmd/api
./bin/api
```

## What Happens on Startup

1. ✅ Database connected
2. 🔄 Running migrations...
   - Creates UUID extension
   - Auto-creates tables (users, courses, enrollments)
   - Creates indexes
   - Sets up foreign keys
3. ✅ Migrations completed
4. 🚀 Server starting on localhost:8080

## Error Handling

GORM errors are now consistent:
- `gorm.ErrRecordNotFound` → mapped to `entity.ErrNotFound`
- All queries use context: `db.WithContext(ctx)`
- Proper error propagation maintained

## Performance Considerations

- Connection pooling configured (5 min, 25 max connections)
- Indexes on frequently queried columns
- Soft deletes for data recovery without affecting queries
- Prepared statements under the hood

## Future Enhancements

Consider adding:
- GORM hooks for audit logging
- Custom validators using GORM callbacks
- Database read replicas
- Caching layer with GORM plugins
- Query optimization with GORM's Debug mode

---

**Migration Completed**: December 15, 2024
**Build Status**: ✅ Successful (14MB binary)
**All Tests**: Ready to implement
