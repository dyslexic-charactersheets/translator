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