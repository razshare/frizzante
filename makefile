configure:
	go run main.go --configure --platform="linux/amd64"

test:
	go run main.go --test

build:
	rm -fr .gen/bin
	mkdir -p .gen/bin
# linux-amd64
	env GOOS=linux GOARCH=amd64 go build -o=.gen/bin/frizzante-linux-amd64 . && \
	zip -9 .gen/bin/frizzante-linux-amd64.zip .gen/bin/frizzante-linux-amd64 && \
	rm .gen/bin/frizzante-linux-amd64
# darwin-amd64
	env GOOS=darwin GOARCH=amd64 go build -o=.gen/bin/frizzante-darwin-amd64 . && \
	zip -9 .gen/bin/frizzante-darwin-amd64.zip .gen/bin/frizzante-darwin-amd64 && \
	rm .gen/bin/frizzante-darwin-amd64
# darwin-arm64
	env GOOS=darwin GOARCH=amd64 go build -o=.gen/bin/frizzante-darwin-arm64 . && \
	zip -9 .gen/bin/frizzante-darwin-arm64.zip .gen/bin/frizzante-darwin-arm64 && \
	rm .gen/bin/frizzante-darwin-arm64

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

publish:
	./publish.sh