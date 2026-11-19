install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1
	echo "install other deps through go tool"

generate-protobufs:
	go tool buf generate

format:
	go tool buf format -w
	gofmt -w -s .

deploy-docker:
	docker compose -f deploy/docker/compose.yaml -f deploy/docker/compose.override.yaml up --build
