APP_EXECUTABLE="./build/out/$(APP)"
APP=new_books_api
APP_VERSION:=0.0.1
APP_COMMIT:=$(shell git rev-parse HEAD)

.PHONY: compile
compile:
	go build -a -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o ./build/out/$(APP) cmd/*.go

.PHONY: deps
deps:
	go mod download

.PHONY: test
test:
	gotestsum --format=testname  --packages ./... --junitfile report.xml -- -coverprofile=coverage.out ./...

.PHONY: clean
clean:
	go clean -testcache
	rm -rf ./build/out

.PHONY: build
build: clean deps test compile

.PHONY: http-serve
http-serve:
	$(APP_EXECUTABLE) -configFile=$(CONFIG_FILE)