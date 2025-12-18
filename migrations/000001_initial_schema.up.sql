-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (
        role IN (
            'admin',
            'instructor',
            'student'
        )
    ),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE
);

-- Create courses table
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    instructor_id UUID NOT NULL,
    difficulty_level VARCHAR(20) NOT NULL CHECK (
        difficulty_level IN (
            'beginner',
            'intermediate',
            'advanced'
        )
    ),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE,
        FOREIGN KEY (instructor_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create modules table
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    course_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE,
        FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
);

-- Create lessons table
CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    module_id UUID NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE,
        FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
);

-- Create lesson_versions table
CREATE TABLE IF NOT EXISTS lesson_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    lesson_id UUID NOT NULL,
    version_number INTEGER NOT NULL,
    content TEXT,
    video_url VARCHAR(500),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (lesson_id) REFERENCES lessons (id) ON DELETE CASCADE,
        CONSTRAINT lesson_version_unique UNIQUE (lesson_id, version_number),
        CONSTRAINT content_or_video CHECK (
            content IS NOT NULL
            OR video_url IS NOT NULL
        )
);

-- Create course_enrollments table
CREATE TABLE IF NOT EXISTS course_enrollments (
    student_id UUID NOT NULL,
    course_id UUID NOT NULL,
    enrollment_date TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        completion_date TIMESTAMP
    WITH
        TIME ZONE,
        status VARCHAR(20) DEFAULT 'active' CHECK (
            status IN (
                'active',
                'completed',
                'dropped'
            )
        ),
        created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE,
        PRIMARY KEY (student_id, course_id),
        FOREIGN KEY (student_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    action VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255),
    resource_type VARCHAR(50),
    user_id UUID,
    payload_before JSONB,
    payload_after JSONB,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);

-- Create indexes for performance

-- User indexes
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email)
WHERE
    deleted_at IS NULL;

-- Course indexes
CREATE INDEX IF NOT EXISTS idx_courses_instructor ON courses (instructor_id)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_courses_difficulty ON courses (difficulty_level)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_courses_created_at ON courses (created_at)
WHERE
    deleted_at IS NULL;

-- Module indexes
CREATE INDEX IF NOT EXISTS idx_modules_course ON modules (course_id, order_index)
WHERE
    deleted_at IS NULL;

-- Lesson indexes
CREATE INDEX IF NOT EXISTS idx_lessons_module ON lessons (module_id)
WHERE
    deleted_at IS NULL;

-- Lesson version indexes
CREATE INDEX IF NOT EXISTS idx_lesson_versions_lesson ON lesson_versions (lesson_id, version_number);

CREATE INDEX IF NOT EXISTS idx_lesson_versions_created ON lesson_versions (created_at);

-- Enrollment indexes
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON course_enrollments (student_id, status)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_enrollments_course ON course_enrollments (course_id, status)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_enrollments_date ON course_enrollments (enrollment_date)
WHERE
    deleted_at IS NULL;

-- Audit log indexes
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs (resource_type, resource_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs (user_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs (created_at);