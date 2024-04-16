install-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.27.2
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

generate-protobufs:
	buf generate

format:
	buf format -w
	gofmt -w -s .

deploy-docker:
	docker compose -f deploy/docker/compose.yaml -f deploy/docker/compose.override.yaml up --build