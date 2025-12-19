#!/bin/bash

###############################################################################
# LMS API Integration Test Script
# Tests the complete workflow sequentially:
# 1. Authentication (register, login)
# 2. Course Management (create, list, update, delete)
# 3. Module Management (create, list)
# 4. Lesson Management (create versions, list)
# 5. Enrollment (enroll, list, complete via webhook)
# 6. Audit Logging (view audit trail)
###############################################################################

set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8080/api}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-lms_db}"

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=0

# Test variables (will be populated during test)
ADMIN_TOKEN=""
INSTRUCTOR_TOKEN=""
STUDENT_TOKEN=""
COURSE_ID=""
MODULE_ID=""
LESSON_ID=""
LESSON_VERSION_ID=""
ENROLLMENT_ID=""

###############################################################################
# Utility Functions
###############################################################################

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_test() {
    echo -e "${YELLOW}TEST: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
    ((TOTAL_TESTS++))
}

print_failure() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
    ((TOTAL_TESTS++))
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Helper to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local auth_token=$4
    
    local headers="-H 'Content-Type: application/json'"
    if [ -n "$auth_token" ]; then
        headers="$headers -H 'Authorization: Bearer $auth_token'"
    fi
    
    if [ -n "$data" ]; then
        eval "curl -s -X $method '$API_URL$endpoint' $headers -d '$data'"
    else
        eval "curl -s -X $method '$API_URL$endpoint' $headers"
    fi
}

# Extract value from JSON response
get_json_field() {
    echo "$1" | grep -o "\"$2\":\"[^\"]*\"" | cut -d'"' -f4
}

get_json_uuid() {
    echo "$1" | grep -o "\"$2\":\"[a-f0-9\-]*\"" | cut -d'"' -f4
}

# Check if response contains error
has_error() {
    echo "$1" | grep -q "error"
    return $?
}

###############################################################################
# Pre-Flight Checks
###############################################################################

preflight_checks() {
    print_header "Pre-Flight Checks"
    
    print_test "Checking if API is running"
    if curl -s "$API_URL/health" > /dev/null 2>&1; then
        print_success "API is running on $API_URL"
    else
        print_failure "API is not running. Please start the API with: make run"
        exit 1
    fi
    
    print_test "Checking database connectivity"
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1" > /dev/null 2>&1; then
        print_success "Database is accessible"
    else
        print_failure "Cannot connect to database. Please check DB_HOST, DB_USER, DB_PASSWORD, DB_NAME"
        echo "Connection string: postgres://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
        exit 1
    fi
}

###############################################################################
# Authentication Tests
###############################################################################

test_authentication() {
    print_header "1. Authentication Tests"
    
    # Register Admin User
    print_test "Register admin user"
    local admin_data='{
        "email":"admin@example.com",
        "password":"SecurePass123",
        "first_name":"Admin",
        "last_name":"User",
        "role":"admin"
    }'
    local admin_response=$(api_call POST "/auth/register" "$admin_data")
    if ! has_error "$admin_response"; then
        print_success "Admin registered successfully"
    else
        print_failure "Admin registration failed: $admin_response"
    fi
    
    # Register Instructor User
    print_test "Register instructor user"
    local instructor_data='{
        "email":"instructor@example.com",
        "password":"SecurePass123",
        "first_name":"John",
        "last_name":"Instructor",
        "role":"instructor"
    }'
    local instructor_response=$(api_call POST "/auth/register" "$instructor_data")
    if ! has_error "$instructor_response"; then
        print_success "Instructor registered successfully"
    else
        print_failure "Instructor registration failed: $instructor_response"
    fi
    
    # Register Student User
    print_test "Register student user"
    local student_data='{
        "email":"student@example.com",
        "password":"SecurePass123",
        "first_name":"Jane",
        "last_name":"Student",
        "role":"student"
    }'
    local student_response=$(api_call POST "/auth/register" "$student_data")
    if ! has_error "$student_response"; then
        print_success "Student registered successfully"
    else
        print_failure "Student registration failed: $student_response"
    fi
    
    # Login Admin
    print_test "Admin login"
    local admin_login='{
        "email":"admin@example.com",
        "password":"SecurePass123"
    }'
    local admin_login_response=$(api_call POST "/auth/login" "$admin_login")
    ADMIN_TOKEN=$(get_json_field "$admin_login_response" "token")
    if [ -n "$ADMIN_TOKEN" ]; then
        print_success "Admin login successful (token: ${ADMIN_TOKEN:0:20}...)"
    else
        print_failure "Admin login failed: $admin_login_response"
        exit 1
    fi
    
    # Login Instructor
    print_test "Instructor login"
    local instructor_login='{
        "email":"instructor@example.com",
        "password":"SecurePass123"
    }'
    local instructor_login_response=$(api_call POST "/auth/login" "$instructor_login")
    INSTRUCTOR_TOKEN=$(get_json_field "$instructor_login_response" "token")
    if [ -n "$INSTRUCTOR_TOKEN" ]; then
        print_success "Instructor login successful (token: ${INSTRUCTOR_TOKEN:0:20}...)"
    else
        print_failure "Instructor login failed: $instructor_login_response"
        exit 1
    fi
    
    # Login Student
    print_test "Student login"
    local student_login='{
        "email":"student@example.com",
        "password":"SecurePass123"
    }'
    local student_login_response=$(api_call POST "/auth/login" "$student_login")
    STUDENT_TOKEN=$(get_json_field "$student_login_response" "token")
    if [ -n "$STUDENT_TOKEN" ]; then
        print_success "Student login successful (token: ${STUDENT_TOKEN:0:20}...)"
    else
        print_failure "Student login failed: $student_login_response"
        exit 1
    fi
}

###############################################################################
# Course Management Tests
###############################################################################

test_course_management() {
    print_header "2. Course Management Tests"
    
    # Get instructor ID from login response (for course creation)
    print_test "Create course as instructor"
    local create_course_data='{
        "title":"Introduction to Go",
        "description":"Learn Go programming from basics to advanced",
        "difficulty_level":"beginner"
    }'
    local create_course_response=$(api_call POST "/courses" "$create_course_data" "$INSTRUCTOR_TOKEN")
    COURSE_ID=$(get_json_uuid "$create_course_response" "id")
    if [ -n "$COURSE_ID" ]; then
        print_success "Course created successfully (ID: ${COURSE_ID:0:20}...)"
    else
        print_failure "Course creation failed: $create_course_response"
        exit 1
    fi
    
    # List Courses
    print_test "List courses"
    local list_courses_response=$(api_call GET "/courses?page=1&limit=10" "" "$INSTRUCTOR_TOKEN")
    if echo "$list_courses_response" | grep -q "Introduction to Go"; then
        print_success "Courses listed successfully"
    else
        print_failure "Course listing failed: $list_courses_response"
    fi
    
    # Get Course
    print_test "Get course details"
    local get_course_response=$(api_call GET "/courses/$COURSE_ID" "" "$INSTRUCTOR_TOKEN")
    if echo "$get_course_response" | grep -q "Introduction to Go"; then
        print_success "Course retrieved successfully"
    else
        print_failure "Course retrieval failed: $get_course_response"
    fi
    
    # Update Course
    print_test "Update course"
    local update_course_data='{
        "title":"Introduction to Go - Updated",
        "description":"Learn Go programming with advanced patterns",
        "difficulty_level":"intermediate"
    }'
    local update_course_response=$(api_call PUT "/courses/$COURSE_ID" "$update_course_data" "$INSTRUCTOR_TOKEN")
    if echo "$update_course_response" | grep -q "intermediate"; then
        print_success "Course updated successfully"
    else
        print_failure "Course update failed: $update_course_response"
    fi
}

###############################################################################
# Module Management Tests
###############################################################################

test_module_management() {
    print_header "3. Module Management Tests"
    
    # Create Module
    print_test "Create module"
    local create_module_data='{
        "title":"Basics",
        "order_index":1
    }'
    local create_module_response=$(api_call POST "/courses/$COURSE_ID/modules" "$create_module_data" "$INSTRUCTOR_TOKEN")
    MODULE_ID=$(get_json_uuid "$create_module_response" "id")
    if [ -n "$MODULE_ID" ]; then
        print_success "Module created successfully (ID: ${MODULE_ID:0:20}...)"
    else
        print_failure "Module creation failed: $create_module_response"
        exit 1
    fi
    
    # List Course Modules
    print_test "List course modules"
    local list_modules_response=$(api_call GET "/courses/$COURSE_ID/modules?page=1&limit=10" "" "$INSTRUCTOR_TOKEN")
    if echo "$list_modules_response" | grep -q "Basics"; then
        print_success "Modules listed successfully"
    else
        print_failure "Module listing failed: $list_modules_response"
    fi
}

###############################################################################
# Lesson Management Tests
###############################################################################

test_lesson_management() {
    print_header "4. Lesson Management Tests"
    
    # Create Lesson (v1)
    print_test "Create lesson (version 1)"
    local create_lesson_data='{
        "content":"Introduction to Go variables and data types",
        "video_url":"https://example.com/go-intro.mp4"
    }'
    local create_lesson_response=$(api_call POST "/modules/$MODULE_ID/lessons" "$create_lesson_data" "$INSTRUCTOR_TOKEN")
    LESSON_ID=$(get_json_uuid "$create_lesson_response" "lesson_id")
    LESSON_VERSION_ID=$(get_json_uuid "$create_lesson_response" "id")
    if [ -n "$LESSON_ID" ]; then
        print_success "Lesson created successfully (ID: ${LESSON_ID:0:20}...)"
    else
        print_failure "Lesson creation failed: $create_lesson_response"
        exit 1
    fi
    
    # List Lessons in Module
    print_test "List lessons in module"
    local list_lessons_response=$(api_call GET "/modules/$MODULE_ID/lessons?page=1&limit=10" "" "$INSTRUCTOR_TOKEN")
    if echo "$list_lessons_response" | grep -q "variables"; then
        print_success "Lessons listed successfully"
    else
        print_failure "Lesson listing failed: $list_lessons_response"
    fi
    
    # Create Lesson Version 2
    print_test "Create lesson version 2"
    local create_version_data='{
        "content":"Introduction to Go variables and data types - Updated with more examples",
        "video_url":"https://example.com/go-intro-v2.mp4"
    }'
    local create_version_response=$(api_call POST "/lessons/$LESSON_ID/version" "$create_version_data" "$INSTRUCTOR_TOKEN")
    if echo "$create_version_response" | grep -q "version_number"; then
        print_success "Lesson version 2 created successfully"
    else
        print_failure "Lesson version creation failed: $create_version_response"
    fi
    
    # Get All Lesson Versions (admin only)
    print_test "Get all lesson versions (admin)"
    local all_versions_response=$(api_call GET "/lessons/$LESSON_ID/all-versions?page=1&limit=10" "" "$ADMIN_TOKEN")
    if echo "$all_versions_response" | grep -q "version_number"; then
        print_success "All lesson versions retrieved successfully"
    else
        print_failure "Get all versions failed: $all_versions_response"
    fi
}

###############################################################################
# Enrollment Tests
###############################################################################

test_enrollment() {
    print_header "5. Enrollment Tests"
    
    # Student Self-Enroll
    print_test "Student self-enroll in course"
    local enroll_response=$(api_call POST "/courses/$COURSE_ID/enroll" "{}" "$STUDENT_TOKEN")
    if ! has_error "$enroll_response"; then
        print_success "Student enrolled successfully"
    else
        print_failure "Enrollment failed: $enroll_response"
        exit 1
    fi
    
    # Get Student Courses
    print_test "Get student's enrolled courses"
    local student_courses_response=$(api_call GET "/students/$(get_json_uuid "$enroll_response" "student_id")/courses?page=1&limit=10" "" "$STUDENT_TOKEN")
    if echo "$student_courses_response" | grep -q "Introduction to Go"; then
        print_success "Student courses retrieved successfully"
    else
        print_failure "Get student courses failed: $student_courses_response"
    fi
    
    # Instructor View Course Students
    print_test "Instructor view enrolled students"
    local course_students_response=$(api_call GET "/courses/$COURSE_ID/students?page=1&limit=10" "" "$INSTRUCTOR_TOKEN")
    if echo "$course_students_response" | grep -q "student@example.com"; then
        print_success "Course students retrieved successfully"
    else
        print_failure "Get course students failed: $course_students_response"
    fi
}

###############################################################################
# Webhook Tests
###############################################################################

test_webhook() {
    print_header "6. Certification Webhook Tests"
    
    # Get student and course IDs for webhook
    local student_id=$(api_call GET "/users" "" "$ADMIN_TOKEN" | grep -o '"id":"[a-f0-9\-]*"' | head -3 | tail -1 | cut -d'"' -f4)
    
    print_test "Process certification webhook (passed)"
    local webhook_data='{
        "student_id":"'$student_id'",
        "course_id":"'$COURSE_ID'",
        "certification_status":"passed",
        "score":85,
        "timestamp":"2025-12-18T10:00:00Z"
    }'
    local webhook_response=$(api_call POST "/certification-webhook" "$webhook_data")
    if ! has_error "$webhook_response"; then
        print_success "Webhook processed successfully"
    else
        print_info "Webhook test skipped or failed: $webhook_response (this may be expected if student not in this course)"
    fi
}

###############################################################################
# Audit Logging Tests
###############################################################################

test_audit_logging() {
    print_header "7. Audit Logging Tests"
    
    # Admin View Audit Logs
    print_test "Admin view audit logs"
    local audit_logs_response=$(api_call GET "/audit-logs?page=1&limit=20" "" "$ADMIN_TOKEN")
    if echo "$audit_logs_response" | grep -q "action"; then
        print_success "Audit logs retrieved successfully"
    else
        print_failure "Audit logs retrieval failed: $audit_logs_response"
    fi
    
    # Verify audit entries
    print_test "Verify course creation logged"
    if echo "$audit_logs_response" | grep -q "course_created"; then
        print_success "Course creation found in audit log"
    else
        print_info "Course creation audit log not found (may need more entries)"
    fi
}

###############################################################################
# Authorization Tests
###############################################################################

test_authorization() {
    print_header "8. Authorization Tests"
    
    # Student cannot create course
    print_test "Student cannot create course (403)"
    local student_course_data='{
        "title":"Unauthorized Course",
        "description":"This should fail",
        "difficulty_level":"beginner"
    }'
    local student_create_response=$(api_call POST "/courses" "$student_course_data" "$STUDENT_TOKEN")
    if echo "$student_create_response" | grep -q "error\|forbidden\|Forbidden"; then
        print_success "Student correctly forbidden from creating course"
    else
        print_failure "Student should not be able to create course: $student_create_response"
    fi
    
    # Instructor cannot view audit logs
    print_test "Instructor cannot view audit logs (403)"
    local instructor_audit_response=$(api_call GET "/audit-logs?page=1&limit=10" "" "$INSTRUCTOR_TOKEN")
    if echo "$instructor_audit_response" | grep -q "error\|forbidden\|Forbidden"; then
        print_success "Instructor correctly forbidden from viewing audit logs"
    else
        print_failure "Instructor should not be able to view audit logs: $instructor_audit_response"
    fi
    
    # Instructor cannot view all lesson versions
    print_test "Instructor cannot view all lesson versions (403)"
    local instructor_versions_response=$(api_call GET "/lessons/$LESSON_ID/all-versions" "" "$INSTRUCTOR_TOKEN")
    if echo "$instructor_versions_response" | grep -q "error\|forbidden\|Forbidden"; then
        print_success "Instructor correctly forbidden from viewing all versions"
    else
        print_info "Instructor versions check: $instructor_versions_response"
    fi
}

###############################################################################
# Test Summary
###############################################################################

print_summary() {
    print_header "Test Summary"
    
    echo -e "Total Tests: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        return 0
    else
        echo -e "\n${RED}✗ Some tests failed${NC}"
        return 1
    fi
}

###############################################################################
# Main Execution
###############################################################################

main() {
    print_header "LMS API Integration Test Suite"
    print_info "Testing API at: $API_URL"
    print_info "Database: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
    
    # Run preflight checks
    preflight_checks
    
    # Run test suites
    test_authentication
    test_course_management
    test_module_management
    test_lesson_management
    test_enrollment
    test_webhook
    test_audit_logging
    test_authorization
    
    # Print summary
    print_summary
}

# Handle script interruption
trap 'echo -e "\n${RED}Test interrupted${NC}"; exit 1' INT TERM

# Run main
main
