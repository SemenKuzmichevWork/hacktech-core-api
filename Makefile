

.PHONY: sqlbuilder
sqlbuilder:
	@echo "Building SQLBuilder"
	@docker compose up -d --wait postgres
	@cd ./migrations && (goose reset || goose up)
	@cd ./migrations && goose up
	@go generate ./internal/storage/...
	@docker compose down