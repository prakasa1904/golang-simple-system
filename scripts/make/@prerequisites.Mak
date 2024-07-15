$(eval WHICH_GAWK := $(strip $(shell which gawk)))
gawk_found = $(strip $(findstring gawk, $(WHICH_GAWK)))

$(eval WHICH_AIR := $(strip $(shell which air)))
air_found = $(strip $(findstring air, $(WHICH_AIR)))

$(eval WHICH_GO := $(strip $(shell which go)))
go_found = $(strip $(findstring go, $(WHICH_GO)))

$(eval WHICH_NODE := $(strip $(shell which node)))
node_found = $(strip $(findstring node, $(WHICH_NODE)))

$(eval WHICH_SQLLITE3 := $(strip $(shell which sqlite3)))
sqllite3_found = $(strip $(findstring sqlite3, $(WHICH_SQLLITE3)))

prerequisite-test: .bin-validator
	@echo "test!"

.bin-validator: ## validate require binaries exist
# exit status code 126 - Command invoked cannot execute
ifneq ($(air_found),air)
	@echo "$(__prefix) install air for better golang devX https://github.com/air-verse/air?tab=readme-ov-file#installation"
	@exit 126
endif

ifneq ($(gawk_found),gawk)
	@echo "$(__prefix) gawk binary not found. install gawk for mac https://formulae.brew.sh/formula/gawk, or install gawk for ubuntu https://howtoinstall.co/en/gawk"
	@exit 126
endif

ifneq ($(go_found),go)
	@echo "$(__prefix) go binary not found, follow golang installation https://go.dev/doc/install"
	@exit 126
endif

ifneq ($(node_found),node)
	@echo "$(__prefix) node binary not found, follow nodeJS installation https://github.com/nvm-sh/nvm?tab=readme-ov-file#installing-and-updating"
	@exit 126
endif

ifneq ($(sqllite3_found),sqlite3)
	@echo "$(__prefix) sqlite3 binary not found, follow sqllite installation https://www.sqlite.org/download.html or using brew for macOS https://formulae.brew.sh/formula/sqlite"
	@exit 126
endif