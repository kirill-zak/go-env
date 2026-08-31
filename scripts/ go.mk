.PHONY: test

lint: # run linter in $dir directory with root config.
	golangci-lint run

test:
	go test ./... -race -cover
