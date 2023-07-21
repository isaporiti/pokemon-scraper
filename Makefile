run:
	go run cmd/main.go

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out

docker_build:
	docker build -t pokemon-scraper .

docker_run:
	docker run -v ./out/:/app/out/ --rm pokemon-scraper