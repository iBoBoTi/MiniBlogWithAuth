run:
	go run cmd/main.go

deploy:
	git push heroku master

migrateup:
	migrate -path db/migration -database "postgresql://postgres:boboti40@localhost:5432/blogar?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:boboti40@localhost:5432/blogar?sslmode=disable" -verbose down