configure:
	go run main.go --configure

test:
	go run main.go --test

build:
	rm -fr .gen/bin
	mkdir -p .gen/bin
# linux-amd64
	cd .gen/bin && env GOOS=linux GOARCH=amd64 go build -o=frizzante-linux-amd64 ../../ && \
	zip -9 frizzante-linux-amd64.zip frizzante-linux-amd64 && \
	rm frizzante-linux-amd64
# linux-arm64
	cd .gen/bin && env GOOS=linux GOARCH=arm64 go build -o=frizzante-linux-arm64 ../../ && \
	zip -9 frizzante-linux-arm64.zip frizzante-linux-arm64 && \
	rm frizzante-linux-arm64
# darwin-amd64
	cd .gen/bin && env GOOS=darwin GOARCH=amd64 go build -o=frizzante-darwin-amd64 ../../ && \
	zip -9 frizzante-darwin-amd64.zip frizzante-darwin-amd64 && \
	rm frizzante-darwin-amd64
# darwin-arm64
	cd .gen/bin && env GOOS=darwin GOARCH=amd64 go build -o=frizzante-darwin-arm64 ../../ && \
	zip -9 frizzante-darwin-arm64.zip frizzante-darwin-arm64 && \
	rm frizzante-darwin-arm64

dev:
	go run main.go --dev

package-watch:
	go run main.go --package-watch

package:
	go run main.go --package

check:
	go run main.go --check

clean:
	go run main.go --clean

format:
	go run main.go --format

install:
	go run main.go --install

update:
	go run main.go --update

hooks:
	go run main.go --hooks