PROJECT_NAME   = github.com/max-grape/test-url-benchmark
CWD            = /go/src/$(PROJECT_NAME)
GOLANG_CI_IMG  = golangci/golangci-lint:v1.42-alpine
GO_IMG         = golang:1.17.1
FROM_IMG       = alpine:3.13.6
IMG            = registry.$(PROJECT_NAME)
TAG            = latest

lint:
	@-docker run --rm -t -w $(CWD) -v $(CURDIR):$(CWD) -e GOFLAGS=-mod=vendor $(GOLANG_CI_IMG) golangci-lint run

unit:
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) \
		$(GO_IMG) sh -c "go list ./... | grep -v 'vendor\|test/' | xargs go test -v"

build: lint unit
	@docker build \
		--build-arg CWD=$(CWD) \
		--build-arg GOOS=linux \
		--build-arg GOARCH=amd64 \
		--build-arg GO_IMG=$(GO_IMG) \
		--build-arg FROM_IMG=$(FROM_IMG) \
		-t $(IMG):$(TAG) .

adown:
	@IMG=$(IMG) TAG=$(TAG) GO_IMG=$(GO_IMG) CWD=$(CWD) docker-compose -f ./test/docker-compose.acceptance.yml down -v --remove-orphans

acceptance: adown
	@IMG=$(IMG) TAG=$(TAG) GO_IMG=$(GO_IMG) FROM_IMG=$(FROM_IMG) CWD=$(CWD) docker-compose -f ./test/docker-compose.acceptance.yml up -d --build --scale acceptance=0
	@IMG=$(IMG) TAG=$(TAG) GO_IMG=$(GO_IMG) FROM_IMG=$(FROM_IMG) CWD=$(CWD) docker-compose -f ./test/docker-compose.acceptance.yml up --abort-on-container-exit acceptance

test: build acceptance

push:
	@docker push $(IMG):$(TAG)

release: build push

deploy:
	#deploy initiated

down:
	@IMG=$(IMG) TAG=$(TAG) GO_IMG=$(GO_IMG) FROM_IMG=$(FROM_IMG) CWD=$(CWD) docker-compose -f ./test/docker-compose.yml down -v --remove-orphans

up: down
	@IMG=$(IMG) TAG=$(TAG) GO_IMG=$(GO_IMG) FROM_IMG=$(FROM_IMG) CWD=$(CWD) docker-compose -f ./test/docker-compose.yml up --build

redis:
	@docker exec -it test_redis_1 redis-cli
