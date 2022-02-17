.PHONY: dockerRun
dockerRun:
	docker run --name=todo-desk -e POSTGRES_PASSWORD='1Asdfghjkl' -p 5436:5432 -d postgres


.PHONY: migrateRun
migrateRun:
	migrate -path ./schema -database 'postgres://postgres:1Asdfghjkl@localhost:5436/postgres?sslmode=disable' up

