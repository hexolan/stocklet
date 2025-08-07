#
# Dev Tools
#

.PHONY: install-tools
install-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.27.2 && \
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0 && \
	npm install

#
# Code Gen
#

.PHONY: gen-protobufs
gen-protobufs:
	buf generate

.PHONY: gen-openapi
gen-openapi: gen-protobufs
	npx openapi-generator-cli generate -g openapi-yaml -i schema/openapi/services.swagger.yaml -o schema/openapi -p outputFile=openapi.yaml

#
# Linting / Formatting
#

.PHONY: format-go
format-go:
	buf format -w && \
	gofmt -w -s .

#
# Deployment
#

.PHONY: deploy-docker
deploy-docker:
	docker compose -f deploy/docker/compose.yaml -f deploy/docker/compose.override.yaml up --build