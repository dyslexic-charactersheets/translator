help:
	@echo "  $ make build                           Compile the app"
	@echo "  $ make run                             Run the app"
	@echo "  $ make test                            Run the app in test mode"

setup:
	go install

build:
	go build src/go/translator.go

run:
	./run-translator.sh

test:
	go run src/go/translator.go

log:
	tailf /var/log/translator

migrate:
	go run src/go/migrate.go

build-all:
	# TODO cross-compile for all targets
	go build src/go/translator.go