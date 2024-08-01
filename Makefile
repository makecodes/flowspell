migration-create:
	@echo "Creating migration file..."
	@migrate create -ext sql -dir db/migrations -seq $(shell echo $(filter-out $@,$(MAKECMDGOALS)))


migration-down:
	@echo "Rolling back migration..."
	@migrate -path db/migrations -database "${DATABASE_URL}" down


migration-up:
	@echo "Applying migration..."
	@migrate -path db/migrations -database "${DATABASE_URL}" up 2

migration-up-test:
	@echo "Applying migration..."
	@migrate -path db/migrations -database "${DATABASE_TEST_URL}" up 2

test:
	@echo "Running tests..."
	@ENV=test go test ./...

coverage:
	@echo "Generating coverage report..."
	@ENV=test go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@echo "Coverage report displayed in terminal"
	@rm coverage.out


.PHONY: migration-create

%:
	@:
