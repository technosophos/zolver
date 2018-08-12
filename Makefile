VERSION := $(shell git describe --tags)

.PHONY: build
build:
	go build -o zolver -ldflags "-X main.version=${VERSION}" server.go

.PHONY: install
install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./zolver ${DESTDIR}/usr/local/bin/zolver

.PHONY: clean
clean:
	rm -f ./zolver.test
	rm -f ./zolver
	rm -rf ./dist

.PHONY: bootstrap
bootstrap:
	# glide install
	go get -u gopkg.in/alecthomas/gometalinter.v1
	# These are failing due to upstream errors.
	#gometalinter.v1 --install
	#dep ensure

.PHONY: test
test:
	go test .

.PHONY: style
style:
	gometalinter.v1 \
		--disable-all \
		--enable deadcode \
		--severity deadcode:error \
		--enable gofmt \
		--enable ineffassign \
		--enable misspell \
		--enable vet \
		--tests \
		--vendor \
  		--deadline 60s \
  		./... || exit_code=1
	gometalinter.v1 \
		--disable-all \
		--enable golint \
		--vendor \
		--skip proto \
		--deadline 60s \
		./... || :

.PHONY: build-release
build-release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/zolver-linux-amd64 -ldflags "-X main.version=${VERSION}" server.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/zolver-windows-amd64 -ldflags "-X main.version=${VERSION}" server.go
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/zolver-darwin-amd64 -ldflags "-X main.version=${VERSION}" server.go
