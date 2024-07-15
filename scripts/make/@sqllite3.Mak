################################################################################
##@ SQLLITE3
################################################################################

sqllite-create: .bin-validator ## Create sqllite database
	@sqlite3 sql_db ".databases"

sqllite-delete: ## Delete sqllite database
	@rm sql_db || printf ""