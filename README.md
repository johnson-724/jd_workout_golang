## Android Play Store

[Link](https://play.google.com/store/apps/details?id=com.govel.workout&hl=zh-TW)

## iOS App Store

[//]: # (emoji angry)

Can't afford the fee. :hankey: :hankey: :hankey:

## Online API document

### Dev

http://www.govel-workout.com:6003/swagger/index.html

### Prod

http://www.govel-workout.com/swagger/index.html

## Sentry

http://www.govel-workout.com:6002/

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