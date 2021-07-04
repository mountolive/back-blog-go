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

todo:
	find . -name '*.go' | xargs grep -n TODO

proto-gen:
	./scripts/proto_all.sh

