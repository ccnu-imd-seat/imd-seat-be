.PHONY: api
api:
	@echo "Generating api..."
	goctl api go -api api/index.api -dir . --style=none

.PHONY: curd
curd:
	@echo "Generating crud"
	goctl model mysql ddl -src internal/sql/data.sql -dir ./internal/model -c --cache=false