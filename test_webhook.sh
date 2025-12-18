#!/bin/bash

# Test script for certification webhook endpoint
# This script demonstrates how to call the webhook with proper HMAC signature

# Configuration
WEBHOOK_URL="http://localhost:8080/api/certification-webhook"
WEBHOOK_SECRET="${WEBHOOK_SECRET:-default-webhook-secret-change-in-production}"

# Sample payload
PAYLOAD='{
  "student_id": "REPLACE_WITH_VALID_STUDENT_UUID",
  "course_id": "REPLACE_WITH_VALID_COURSE_UUID",
  "certification_status": "passed",
  "score": 95,
  "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'

echo "Payload:"
echo "$PAYLOAD" | jq .

# Generate HMAC signature
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | awk '{print $2}')

echo ""
echo "HMAC Signature: $SIGNATURE"
echo ""

# Send webhook request
echo "Sending webhook request..."
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD" \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "Instructions:"
echo "1. Start the server: go run ./src/cmd/api"
echo "2. Create a student user and enroll them in a course"
echo "3. Update this script with valid student_id and course_id"
echo "4. Run this script: ./test_webhook.sh"
echo ""
echo "Environment variables:"
echo "  WEBHOOK_SECRET - The shared secret for HMAC validation (default: default-webhook-secret-change-in-production)"
