## Online API document

http://52.192.216.234:6003/swagger/index.html

## Run docker

* `docker-compose build`
* `docker-compose up -d`

## Hot reload

use command `air`

### test api Hot reload

run `air -c .air_test.toml`

## Run test api

run `go run cmd/test-api.go` in the container.

## Migrate

### create new migration

`migrate create -ext sql -dir database/migrations create_users_table`

### run migrate
- up
`migrate -database "mysql://${username}:${password}@tcp${MYSQL_URL}/${DB_NAME}" -path db/migrations up`


- down
`migrate -database "mysql://${username}:${password}@tcp${MYSQL_URL}/${DB_NAME}" -path db/migrations down ${BATCH}`

## Api document

### Generate document

`swag init -g cmd/main.go`

### Check in browser

http://localhost:8080/swagger/index.html