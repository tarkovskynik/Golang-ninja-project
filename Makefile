lint:
	staticcheck ./...

test:
	# go test -timeout 30s
	# go test -v -cover ./...
	go test ./... -coverprofile cover_file
	go test --short -coverprofile=cover_file -v ./...
	del cover_file

test.coverage:
	go test ./... -coverprofile cover_file
	go tool cover -func=cover_file
	del cover_file

swag:
	swag init -g cmd/main.go


build:
	docker-compose build gnp

run:
	docker-compose up gnp
