frizzante:
	cd template/project && make frizzante

configure:
	cd template/project && make configure

test:
	go run main.go -y --configure --app="template/project/app"
	make package
	test -d template/project/app && \
	rm -fr app && \
	rm -fr svelte/ssr/app && \
	mkdir -p app && \
	mkdir -p svelte/ssr/app && \
	cp -r template/project/app .
	cp -r template/project/app svelte/ssr
	go run main.go -y --test

dev:
	cd template/project && make dev

package-watch:
	cd template/project && make package-watch

package:
	cd template/project && make package

check:
	go run main.go --check

clean:
	cd template/project && make clean

format:
	go run main.go --format

install:
	cd template/project && make install

update:
	cd template/project && make update

hooks:
	go run main.go --hooks

zip:
	./zip.sh

publish: zip
	./publish.sh