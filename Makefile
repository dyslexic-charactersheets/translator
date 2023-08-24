help:
	@echo "  $ make build                           Compile the app"
	@echo "  $ make run                             Run the app"
	@echo "  $ make test                            Run the app in test mode"

setup:
	go install

build:
	go build translator.go

run:
	./run-translator.sh

test:
	go run translator.go

log:
	tailf /var/log/translator

migrate:
	go run src/go/migrate.go

