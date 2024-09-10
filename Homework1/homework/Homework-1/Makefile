ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=1.6.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: run
run:
	go run ./cmd/main.go -service

.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif


.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

.PHONY: gofmt
gofmt:
	goimports -l -w $(CURDIR)

#.PHONY: test
##	go test ./...


.PHONY: test
test:

	cd ./internal/infrastructure/kafka && go test


build:
	docker-compose build

up-all:
	docker-compose up -d zookeeper kafka1 kafka2 kafka3

down:
	docker-compose down


#cd ./internal/pkg/service && go test
#cd ./internal/pkg/repository/postgresql && go test
#cd ./internal/storage && go test