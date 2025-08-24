test:
	rm -fr app
	rm -fr svelte/ssr/app
	cp -r template/project/app .
	cp -r template/project/app svelte/ssr
	go test . && \
	go test ./send && \
	go test ./receive && \
	go test ./svelte/ssr

configure:
	go run main.go --app="template/project/app" --configure
	rm -fr template/project/.gen
	cp -r .gen template/project

dev:
	go run main.go --app="template/project/app" --dev

package-watch:
	go run main.go --app="template/project/app" --package-watch

package:
	go run main.go --app="template/project/app" --package

check:
	go run main.go --app="template/project/app" --check

clean:
	go run main.go --app="template/project/app" --clean-project

format:
	go run main.go --app="template/project/app" --format

install:
	go run main.go --app="template/project/app" --install

update:
	go run main.go --app="template/project/app" --update

zip:
	./zip.sh

publish: zip
	./publish.sh