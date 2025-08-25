test:
	rm -fr svelte/ssr/app
	cp -r template/project/app svelte/ssr
	cp -r template/project/app js
	go run main.go --test

configure:
	go run main.go --app="template/project/app" --configure
	rm -fr template/project/.gen
	cp -r .gen template/project

dev:
	cd template/project && make dev

package-watch:
	cd template/project && make package-watch

package:
	cd template/project && make package

check:
	cd template/project && make check

build:
	cd template/project && make build

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