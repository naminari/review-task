# Makefile для PR Reviewer Service

.PHONY: help build run test unit-test integration-test quick-test examples clean docker-up docker-down docker-logs health status

build:
	go build -o bin/review-service main.go

run:
	go run main.go

test: unit-test quick-test

unit-test:
	go test -v ./...

quick-test:
	chmod +x scripts/quick_test.sh
	./scripts/quick_test.sh

examples:
	chmod +x scripts/clean_examples.sh
	./scripts/clean_examples.sh

clean-examples:
	@echo "Очистка тестовых данных..."
	@docker-compose exec db psql -U postgres -d review_service -c "DELETE FROM pull_requests WHERE pull_request_id LIKE 'pr-%';" || true
	@docker-compose exec db psql -U postgres -d review_service -c "DELETE FROM users WHERE user_id LIKE 'user%-%';" || true
	@docker-compose exec db psql -U postgres -d review_service -c "DELETE FROM teams WHERE team_name LIKE 'team-%';" || true

clean:
	rm -rf bin/

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

health:
	curl -s http://localhost:8080/health | jq .

status:
	docker-compose ps