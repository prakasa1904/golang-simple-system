################################################################################
##@ BUILD
################################################################################

build-app: ## Build service for production
	@go build -o app-boilerplate cmd/webapp/*.go

build-css: ## Build css for production
	@npm run build
