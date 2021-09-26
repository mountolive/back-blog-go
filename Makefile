test:
	./scripts/test_all.sh

test-post:
	cd post && make test

test-user:
	cd user && make test

test-gateway:
	cd gateway && cargo test

clippy:
	cd gateway && cargo clippy

build-gateway:
	cd gateway && cargo build

todo:
	find . -name '*.go' -or -name '*.rs' | xargs grep -n TODO

proto-gen:
	./scripts/proto_all.sh

proto-compile:
	./scripts/proto_compile.sh

tidy:
	cd user && go mod tidy
	cd ..
	cd post && go mod tidy

start-all:
	docker-compose up -d

rebuild:
	docker-compose down -v
	docker-compose rm
	docker-compose build

start-all-block:
	docker-compose up

shutdown-all:
	docker-compose down --remove-orphans
