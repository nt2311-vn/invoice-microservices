$APP_NAME = "cronjob-invoice"

build:
	@echo "Building the binary..."
	@go build -o bin/$(APP_NAME)
