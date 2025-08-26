APP_NAME := $(shell grep '^module' go.mod | cut -d' ' -f2 | tr -d '"')
PID := $(shell ps -ef | grep $(APP_NAME) | grep -v grep | awk '{print $$2}')

install:
	@echo "Install dependencies..."
	@go mod download

build:
	@echo "Building application... ${APP_NAME}"
	@go build -o ./build/${APP_NAME} .

run-prod:
	@echo "Running application... ${APP_NAME}"
	@mkdir -p logs
	@nohup ./build/$(APP_NAME) >> ./logs/${APP_NAME}.log 2>&1 &

run-dev:
	@echo "Running application in development mode... ${APP_NAME}"
	@go run main.go

run-watch:
	@echo "Running application with file watching... ${APP_NAME}"
	@air

stop:
	@echo "Stopping application... ${APP_NAME} | PID: $(PID)"
	@if [ -n "$(PID)" ]; then \
		kill -9 $(PID); \
		echo "Application stopped."; \
	else \
		echo "No running application found."; \
	fi
	@make clean

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build
	@rm -rf tmp

restart: stop clean build run-prod