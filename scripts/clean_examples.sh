#!/bin/bash

echo "Рабочие примеры использования PR Reviewer API"
echo "============================================="
echo "Используем существующие тестовые данные"
echo

BASE_URL="http://localhost:8080"
TIMESTAMP=$(date +%s)

PR_ID="pr-example-${TIMESTAMP}"

echo "Используемые данные:"
echo "- Команда: backend (u1, u2, u3, u4)"
echo "- PR: $PR_ID"
echo

echo "1. Получаем информацию о команде backend"
echo "-----------------------------------------"
curl -s "$BASE_URL/team/get?team_name=backend" | jq '.'
echo

echo "2. Создаем Pull Request (автоназначение ревьюеров)"
echo "--------------------------------------------------"
curl -X POST "$BASE_URL/pullRequest/create" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "'"$PR_ID"'",
    "pull_request_name": "Example Feature",
    "author_id": "u1"
  }'
echo
echo

sleep 1

echo "3. Смотрим PR пользователя u2"
echo "-----------------------------"
curl -s "$BASE_URL/users/getReview?user_id=u2" | jq '.'
echo

sleep 1

echo "4. Переназначаем ревьюера"
echo "--------------------------"
curl -X POST "$BASE_URL/pullRequest/reassign" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "'"$PR_ID"'", 
    "old_user_id": "u2"
  }'
echo
echo

sleep 1

echo "5. Смотрим обновленный PR"
echo "--------------------------"
curl -s "$BASE_URL/users/getReview?user_id=u3" | jq '.'
echo

sleep 1

echo "6. Мерджим PR"
echo "-------------"
curl -X POST "$BASE_URL/pullRequest/merge" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "'"$PR_ID"'"
  }'
echo
echo

sleep 1

echo "7. Пытаемся изменить мерджнутый PR (должна быть ошибка)"
echo "-------------------------------------------------------"
curl -X POST "$BASE_URL/pullRequest/reassign" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "'"$PR_ID"'",
    "old_user_id": "u3"
  }'
echo
echo

sleep 1

echo "8. Деактивируем пользователя"
echo "----------------------------"
curl -X POST "$BASE_URL/users/setIsActive" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u4",
    "is_active": false
  }'
echo
echo

sleep 1

echo "9. Проверяем обновленного пользователя"
echo "--------------------------------------"
curl -s "$BASE_URL/team/get?team_name=backend" | jq '.members[] | select(.user_id == "u4")'
echo

echo "Примеры завершены!"
echo "Созданный PR: $PR_ID"