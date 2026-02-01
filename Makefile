.SILENT: create_db psql migrate_up_all migrate_down_all create_seed run_seeds

include .env
export $(shell sed 's/=.*//' .env)

MIGRATION_DIR="tools/migrations"
SEEDS_DIR="tools/seeds"
TMP_COVER_FILE="./tmp/coverage.txt"
TMP_COVER_HTML_FILE="./tmp/index.html"

migrate_bin:
	@echo -e "\n\n\n======================="
	@echo -e 	   "Building migration binary"
	@echo -e 	   "=======================\n"
	@go build -o bin/migrate_bin ./cmd/main.go 
	@echo -e "\n\n\n======================="
	@echo -e 	   "Migration binary built successfully"
	@echo -e 	   "=======================\n"
	@./bin/migrate_bin job migrate up tools/migrations

up:
	@docker compose -f ./tools/development/docker-compose_infrastructure.yml up -d
	@docker compose -f ./tools/development/docker-compose_app.yml up -d

down:
	@docker compose -f ./tools/development/docker-compose_infrastructure.yml down
	@docker compose -f ./tools/development/docker-compose_app.yml down

clean-data:
	@rm -rf tools/development/.cache
	@rm -rf tools/development/.db
	@rm -rf tools/development/.emqx

clean-image:
	@docker image rm -f development-migrate
	@docker image rm -f development-api
	@docker image rm -f development-kafka_consumer
	@docker image rm -f development-mqtt_subscriber

infrastructure:
	@docker compose -f ./tools/development/docker-compose_infrastructure.yml up -d

command:
	@./tools/script/run.sh

http:
	@go run ./cmd/main.go http echo

tests:
	@echo -e "\n\n\n======================="
	@echo -e 	   "Running full test suite"
	@echo -e 	   "=======================\n"
	@go test -timeout 2m -cover ./... | sed ''/^ok/s//$$(printf "\033[32mok\033[0m")/'' | sed ''/^?/s//$$(printf "\033[33m-\033[0m")/'' | sed ''/^FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

mocks:
	@go generate ./...

coverage:
	@echo -e "\n==============================================="
	@echo -e   "Evaluating coverage for $$(printf "\033[33mpackage as part of repo\033[0m")"
	@echo -e   "===============================================\n"
	
	@go test -timeout 5m -coverpkg=./... -coverprofile=$(TMP_COVER_FILE) ./... | sed ''/^ok/s//$$(printf "\033[32mok\033[0m")/'' | sed ''/^?/s//$$(printf "\033[33m-\033[0m")/'' | sed ''/^FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
	@cat $(TMP_COVER_FILE) | grep -v mock > coverage.final.txt
	@mv coverage.final.txt $(TMP_COVER_FILE)

	@echo -e "\n====================="
	@echo -e   "$$(printf "\033[33mBreakdown by function\033[0m")"
	@echo -e   "=====================\n"
	@go tool cover -func=$(TMP_COVER_FILE)
	@go tool cover -o $(TMP_COVER_HTML_FILE) -html=$(TMP_COVER_FILE)

	@echo -e "\nFor a full visual breakdown, see $$(printf "\033[32mhttp://localhost:8082\033[0m") for more details"

	@echo -e "\n"

lint:
	@./bin/golangci-lint run
	

migration:
ifndef name
	$(error name is not set)
else
	migrate create -ext sql -dir "$(MIGRATION_DIR)" $(name)
endif

migrate_up:
	driver = ${MASTER_DATABASE_DRIVER}
	if [driver -eq 'nrpostgres']
	then
		driver = 'postgres'
	fi
ifndef step
	migrate -source file://$(MIGRATION_DIR) -database "${driver}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" up 1
else
	migrate -source file://$(MIGRATION_DIR) -database "${MASTER_DATABASE_DRIVER}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" up $(step)
endif

migrate_down: 
ifndef step
	migrate -source file://$(MIGRATION_DIR) -database "${MASTER_DATABASE_DRIVER}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" down 1
else
	migrate -source file://$(MIGRATION_DIR) -database "${MASTER_DATABASE_DRIVER}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" down $(step)
endif

migrate_up_all:
	migrate -source file://$(MIGRATION_DIR) -database "${MASTER_DATABASE_DRIVER}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" up

migrate_down_all:
	migrate -source file://$(MIGRATION_DIR) -database "${MASTER_DATABASE_DRIVER}://${MASTER_DATABASE_USERNAME}:${MASTER_DATABASE_PASSWORD}@${MASTER_DATABASE_HOST}:${MASTER_DATABASE_PORT}/${MASTER_DATABASE_NAME}?sslmode=${MASTER_DATABASE_SSL}" down

create_seed:
ifndef name
	$(error name is not set)
else
	touch "$(SEEDS_DIR)/$$(date +"%Y%m%d%H%M%S")_$(name).sql"
endif

run_seeds:
	for file in $(SEEDS_DIR)/*; do PGPASSWORD=${MASTER_DATABASE_PASSWORD} psql -h ${MASTER_DATABASE_HOST} -p ${MASTER_DATABASE_PORT} -U ${MASTER_DATABASE_USERNAME} ${MASTER_DATABASE_NAME} -f "$$file"; done

compose-build:
	docker compose build --no-cache --progress plain

compose-up:
	docker compose up -d

compose-down:
	docker compose down

clean: 
	rm -rf cmd internal tmp go.*