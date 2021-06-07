proto-gen:
	./scripts/proto-gen.sh	
test:
	go test --race ./...
mocks:
	# TODO Add mocks to all packages and switch old mock implementations
	moq -out httpx/mock_repository_test.go -pkg httpx_test usecase Repository
