################################################################################
##@ PLAYGROUND
################################################################################

init: .bin-validator ## Init development runtime
	@touch panji_express.db
	@go mod tidy
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

clean-up: .bin-validator ## Cleanup database, create new with make init
	@rm panji_express.db
	@rm panji_wa.db
