## Run test api

docker-compose up -d

run `go run test-api.go` in the container.

## Migrate

### create new migration

`migrate create -ext sql -dir database/migrations create_users_table`

### run migrate
`migrate -database ${MYSQL_URL} -path db/migrations up`
