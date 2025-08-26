test:
	go run main.go --test

sync: 
	cp -r template/project/app svelte/ssr
	cp -r template/project/app js
	cp -r .gen template/project

configure:
	go run main.go --app="template/project/app" --configure

dev:
	cd template/project && make dev

package-watch:
	go run main.go --app="template/project/app" --package-watch

package:
	go run main.go --app="template/project/app" --package

check:
	go run main.go --app="template/project/app" --check

clean:
	rm -fr js/app
	rm -fr svelte/ssr/app
	rm -fr template/project/.gen
	go run main.go --app="template/project/app" --clean-project

format:
	go run main.go --app="template/project/app" --format

install:
	go run main.go --app="template/project/app" --install
	cd template/project && make install

update:
	go run main.go --app="template/project/app" --update
	cd template/project && make update

coverage:
	go test ./... -coverprofile cover.out && \
	go tool cover -html=cover.out -o cover.html