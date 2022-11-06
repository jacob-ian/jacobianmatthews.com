migration_create:
	migrate create -ext sql -dir postgres/migrations -seq $(NAME)
migration_up:
	migrate -database $(DB_CONNECTION_STRING) -path postgres/migrations up
migration_down:
	migrate -database $(DB_CONNECTION_STRING) -path postgres/migrations down
