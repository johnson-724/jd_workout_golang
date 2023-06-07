test:
	go test -cover ./... --coverprofile=coverage.txt -race -covermode=atomic -cover=true