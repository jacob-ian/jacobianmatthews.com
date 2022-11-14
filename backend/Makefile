migration_create:
	migrate create -ext sql -dir internal/postgres/migrations -seq $(NAME)
migration_up:
	migrate -database $(DB_CONNECTION_STRING) -path internal/postgres/migrations up
migration_down:
	migrate -database $(DB_CONNECTION_STRING) -path internal/postgres/migrations down
migration_force:
	migrate -database $(DB_CONNECTION_STRING) -path internal/postgres/migrations force $(VERSION)
seed_dev:
	./scripts/seed_db.sh $(DB_CONNECTION_STRING) internal/postgres/seed/dev

