compile:
	go build -v ./...

unit_tests:
	go clean -testcache
	docker-compose -f cfg/docker-compose.yml up -d
	export MONGODB_URL="mongodb://admin:password@localhost:27017" && go test -v ./...
	docker-compose -f cfg/docker-compose.yml stop
