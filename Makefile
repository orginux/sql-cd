build-go:
	@go build -o bin/sql-cd main.go

build-image:
	@docker build --tag sql-cd:latest .

# Apply queries
test-apply: build-go
	./bin/sql-cd \
		-host "${HOST:-localhost}" \
		-git-url "${REPO}"
		-git-path "${PATH}"
