#!/bin/bash

echo "🚀 Быстрый тест API"
echo "==================="
echo "Используем существующие тестовые данные:"
echo "- Команда: backend (u1, u2, u3, u4)"
echo "- Команда: frontend (u5, u6, u7)"
echo

BASE_URL="http://localhost:8080"

echo "1. Health check..."
curl -s "$BASE_URL/health" | jq '.' || echo "Health check failed"
echo

echo "2. Получаем команду backend..."
curl -s "$BASE_URL/team/get?team_name=backend" | jq '.' || echo "Failed to get team"
echo

echo "3. Создаем PR (автор u1 из команды backend)..."
curl -s -X POST "$BASE_URL/pullRequest/create" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-quick-001",
    "pull_request_name": "Quick Test Feature", 
    "author_id": "u1"
  }' | jq '.' || echo "Failed to create PR"
echo

echo "4. PR пользователя u2..."
curl -s "$BASE_URL/users/getReview?user_id=u2" | jq '.' || echo "Failed to get user PRs"
echo

echo "5. Мерджим PR..."
curl -s -X POST "$BASE_URL/pullRequest/merge" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-quick-001"
  }' | jq '.' || echo "Failed to merge PR"
echo

echo "6. Пытаемся изменить мерджнутый PR..."
curl -s -X POST "$BASE_URL/pullRequest/reassign" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-quick-001",
    "old_user_id": "u2"
  }' | jq '.' || echo "Expected error for merged PR"
echo

echo "✅ Быстрый тест завершен"
