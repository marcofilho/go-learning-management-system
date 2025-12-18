-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_created;

DROP INDEX IF EXISTS idx_audit_logs_user;

DROP INDEX IF EXISTS idx_audit_logs_resource;

DROP INDEX IF EXISTS idx_audit_logs_action;

DROP INDEX IF EXISTS idx_enrollments_date;

DROP INDEX IF EXISTS idx_enrollments_course;

DROP INDEX IF EXISTS idx_enrollments_student;

DROP INDEX IF EXISTS idx_lesson_versions_created;

DROP INDEX IF EXISTS idx_lesson_versions_lesson;

DROP INDEX IF EXISTS idx_lessons_module;

DROP INDEX IF EXISTS idx_modules_course;

DROP INDEX IF EXISTS idx_courses_created_at;

DROP INDEX IF EXISTS idx_courses_difficulty;

DROP INDEX IF EXISTS idx_courses_instructor;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_is_active;

DROP INDEX IF EXISTS idx_users_role;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS audit_logs;

DROP TABLE IF EXISTS course_enrollments;

DROP TABLE IF EXISTS lesson_versions;

DROP TABLE IF EXISTS lessons;

DROP TABLE IF EXISTS modules;

DROP TABLE IF EXISTS courses;

DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";