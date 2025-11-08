#!/bin/bash

# 회원가입 테스트
echo "=== Testing Signup ==="
curl -X POST http://localhost:8080/account/signup \
  -H "Content-Type: application/cbor" \
  -H "Accept: application/cbor" \
  -d '{"email": "test@example.com", "password": "password123", "nickname": "TestUser"}' \
  -v

echo -e "\n\n=== Testing Login ==="
# 로그인 테스트
curl -X POST http://localhost:8080/account/login \
  -H "Content-Type: application/cbor" \
  -H "Accept: application/cbor" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt \
  -v

echo -e "\n\n=== Testing Me ==="
# 인증된 사용자 정보 확인
curl -X GET http://localhost:8080/account/me \
  -H "Accept: application/cbor" \
  -b cookies.txt \
  -v
