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
	glide install

.PHONY: test
test:
	go test .

.PHONY: build-release
build-release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/zolver-linux-amd64 -ldflags "-X main.version=${VERSION}" server.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/zolver-windows-amd64 -ldflags "-X main.version=${VERSION}" server.go
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/zolver-darwon-amd64 -ldflags "-X main.version=${VERSION}" server.go
