#!/bin/bash

# Test script for refresh token endpoint
# This script tests the refresh token functionality manually

BASE_URL="http://localhost:3000"

echo "=== Testing Refresh Token Endpoint ==="
echo ""

# Step 1: Register a test user
echo "1. Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')
echo "Register response: $REGISTER_RESPONSE"
echo ""

# Step 2: Login to get tokens
echo "2. Logging in to get initial tokens..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')
echo "Login response: $LOGIN_RESPONSE"
echo ""

# Extract refresh token from response
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)
echo "Extracted refresh token: $REFRESH_TOKEN"
echo ""

# Step 3: Use refresh token to get new tokens
echo "3. Using refresh token to get new tokens..."
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
echo "Refresh response: $REFRESH_RESPONSE"
echo ""

# Step 4: Try to use the old refresh token again (should fail)
echo "4. Trying to reuse old refresh token (should fail due to token rotation)..."
REUSE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
echo "Reuse response: $REUSE_RESPONSE"
echo ""

# Step 5: Test with invalid token
echo "5. Testing with invalid refresh token..."
INVALID_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"invalid-token-12345"}')
echo "Invalid token response: $INVALID_RESPONSE"
echo ""

echo "=== Test completed ==="
