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

run-users:
	cd user && go run ./...

run-posts:
	cd post && go run ./...

run-gateway:
	make build-gateway
	cd gateway && cargo run

local-infra:
	docker-compose --file docker-compose-infra.yml --env-file .env.local up -d

local-infra-no-d:
	docker-compose --file docker-compose-infra.yml --env-file .env.local up

down-infra:
	docker-compose --file docker-compose-infra.yml down --remove-orphans

restart-infra:
	make down-infra
	make local-infra

ps-infra:
	docker-compose --file docker-compose-infra.yml ps

logs-infra:
	docker-compose --file docker-compose-infra.yml logs -f

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

build-posts:
	cd post && $(MAKE) build

build-users:
	cd user && $(MAKE) build

build-gateway-release:
	cd gateway && $(MAKE) build-release
