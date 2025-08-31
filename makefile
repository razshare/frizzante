test:
	go test ./cli/app/...
	go test ./cli/extension/...
	go test ./cli/menu/...
	go test ./cli/npm/...
	go test ./cli/path/...
	go test ./cli/user/...
	go test ./embeds/...
	go test ./files/...
	go test ./js/...
	go test ./js/runtime/...
	go test ./mime/...
	go test ./receive/...
	go test ./send/...
	go test ./server/...
	go test ./svelte/ssr/...
	go test ./text/...

lock:
	mv internal/template/project/go.mod.txt internal/template/project/go.mod &
	mv internal/template/project/go.sum.txt internal/template/project/go.sum

unlock:
	mv internal/template/project/go.mod internal/template/project/go.mod.txt &
	mv internal/template/project/go.sum internal/template/project/go.sum.txt

sync: 
	cp -r internal/template/project/app svelte/ssr
	cp -r internal/template/project/app js
	cp -r .gen internal/template/project

configure:
	go run main.go --app="internal/template/project/app" --configure

dev:
	cd internal/template/project && make dev

package-watch:
	go run main.go --app="internal/template/project/app" --package-watch

package:
	go run main.go --app="internal/template/project/app" --package

check:
	go run main.go --app="internal/template/project/app" --check

clean:
	rm -fr js/app
	rm -fr svelte/ssr/app
	rm -fr internal/template/project/.gen
	go run main.go --app="internal/template/project/app" --clean-project

format:
	go run main.go --app="internal/template/project/app" --format

install:
	go run main.go --app="internal/template/project/app" --install
	cd internal/template/project && make install

update:
	go run main.go --app="internal/template/project/app" --update
	cd internal/template/project && make update

reset:
	go run main.go --reset

coverage:
	go test ./... -coverprofile cover.out && \
	go tool cover -html=cover.out -o cover.html