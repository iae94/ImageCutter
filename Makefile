UNIT_TEST_DIR=pkg\lru
INTEGRATION_TEST_DIR=pkg\integration_tests
CUTTER_DIR=cmd\cutter
DOCKER_DIR=docker

all: build test run
test: unit_test integration_test
unit_test:
		@echo "Run unit tests(lru)..."
		@cd $(UNIT_TEST_DIR) && \
		go test -v
integration_test:
		@echo "Run integration tests..."
		@cd $(INTEGRATION_TEST_DIR)
		docker-compose -f $(DOCKER_DIR)/docker-compose.test.yaml up tests
		docker-compose -f $(DOCKER_DIR)/docker-compose.test.yaml down

build:
		@echo "Build image cutter service..."
		@echo "Check go vet..."
		@cd $(CUTTER_DIR) && go vet -v
		@echo "Check golangci-lint..."
		@cd $(CUTTER_DIR) && golangci-lint run -v
		@echo "Build with -race..."
		@cd $(CUTTER_DIR) && go build -race -v .
run:
		@echo "Start cutter service with docker-compose..."
		docker-compose -f $(DOCKER_DIR)\docker-compose.yaml up cutter
stop:
		@echo "Down cutter service with docker-compose..."
		docker-compose -f $(DOCKER_DIR)\docker-compose.yaml stop cutter
