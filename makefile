test:
	rm -fr js/app
	rm -fr svelte/ssr/app
	cp -r template/project/app svelte/ssr
	cp -r template/project/app js
	go run main.go --test

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

coverage:
	go test ./... -coverprofile cover.out && \
	go tool cover -html=cover.out -o cover.html

zip:
	./zip.sh

publish: zip
	./publish.sh