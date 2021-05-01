test:
	./scripts/test_all.sh

todo:
	find . -name '*.go' | xargs grep -n TODO
