default: vet build

all: vet insider

build: insider

fmt:
	find . ! -path "*/vendor/*" -type f -name '*.go' -exec gofmt -l -s -w {} \;

vet:
	go vet ./...

insider:
	go build ${GCFLAGS} -ldflags "${LDFLAGS}" ./cmd/insider

clean:
	rm -vf ./insider

format:
	@echo "Formatinsiderg Go files..."
	@find . -name "*.go" -type f -print0 | xargs -0 gofmt -s -w
	@goimports -w .

swaggen:
	swag fmt && swag init --parseDependency -g ./cmd/insider/main.go -o ./api/docs

test:
	go test -race ./...

coverage:
	go test -coverprofile=c.out -race ./...
	go tool cover -html="c.out"

.NOTPARALLEL:
