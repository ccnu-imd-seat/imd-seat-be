.PHONY: api
api:
	@echo "Generating api..."
	goctl api go -api api/index.api -dir . --style=none

.PHONY: curd
struct:
	@echo "Generating crud"
	goctl model mysql ddl -src sql/data.sql -dir ./internal/model -c