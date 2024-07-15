################################################################################
##@ PLAYGROUND
################################################################################

init: .bin-validator ## Init development runtime
	@touch sql_db
	@sleep 5
	@make migrate

run: .bin-validator ## Run service development
	@export ENV=development && go run cmd/webapp/*.go

css-dev: ## Run css development
	@npm run watch

run-hot: .bin-validator ## Run service development with hot reload
	@air

migrate: .bin-validator ## Migrate database "make migrate env=beta"
	@export ENV=$(env) && go run cmd/migrate/*.go