.PHONY: api
api:
	@echo "Generating api..."
	goctl api go --api api/index.api --dir . 

.PHONY: curd
curd:
	@echo "Generating crud"
	goctl model mysql ddl -src internal/sql/data.sql -dir ./internal/model -c --cache=false

.PHONY: reapply
reapply:
	@echo "Reapplying ..."
	docker compose down
	docker rmi imd-seat-be-app
	CGO_ENABLED=0 go build -o imd-be main.go
	docker compose up -d