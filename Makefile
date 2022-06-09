lint:
	@echo "Running linter checks"
	golangci-lint run

test:
	@echo "Running UNIT tests"
	@go clean -testcache
	go test -cover -race -short ./... | { grep -v 'no test files'; true; }

cover-html:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -html=coverage.out

cover:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -func=coverage.out