test:
	go run main.go --test

configure:
	cd template/project && make configure

dev:
	cd template/project && make dev

package-watch:
	cd template/project && make package-watch

package:
	cd template/project && make package

copies:
	rm -fr .gen && \
	rm -fr app && \
	rm -fr svelte/ssr/app && \
	cp -r template/project/.gen .
	cp -r template/project/app .
	cp -r template/project/app svelte/ssr

check:
	cd template/project && make check

clean:
	cd template/project && make clean

format:
	cd template/project && make format

install:
	cd template/project && make install

update:
	cd template/project && make update

zip:
	./zip.sh

publish: zip
	./publish.sh