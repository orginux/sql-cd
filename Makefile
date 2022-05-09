build-go:
	@go build -o bin/sql-cd main.go

build-image:
	@docker build --tag sql-cd:latest .

# Test enviroment
test-server-up:
	@cd tests/ && docker-compose up -d --build && docker-compose ps

test-server-destroy:
	@cd tests/ && docker-compose down


# Apply queries
test-apply: build-go
	./bin/sql-cd \
		-host "${HOST:-localhost}" \
		-git-url "${REPO}"
		-git-path "${PATH}"
