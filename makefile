test:
	go run main.go --test

configure:
	TRACE=1 cd template/project && make configure

dev:
	TRACE=1 cd template/project && make dev

package-watch:
	TRACE=1 cd template/project && make package-watch

package:
	TRACE=1 cd template/project && make package

copies:
	rm -fr .gen && \
	rm -fr app && \
	rm -fr svelte/ssr/app && \
	cp -r template/project/.gen .
	cp -r template/project/app .
	cp -r template/project/app svelte/ssr

check:
	TRACE=1 cd template/project && make check

clean:
	TRACE=1 cd template/project && make clean

format:
	TRACE=1 cd template/project && make format

install:
	TRACE=1 cd template/project && make install

update:
	TRACE=1 cd template/project && make update

zip:
	./zip.sh

publish: zip
	./publish.sh