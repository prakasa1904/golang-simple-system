package group

var (
	SelectColumn       = "id, name, status, created_at, updated_at"
	AllowedFilterQuery = []string{"id", "name"}
)

var (
	StatusDeleted   = 0
	StatusActivated = 1
)
