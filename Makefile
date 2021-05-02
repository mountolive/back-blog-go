test:
	./scripts/test_all.sh

todo:
	find . -name '*.go' | xargs grep -n TODO

proto-gen:
	./scripts/proto_all.sh

