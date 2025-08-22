frizzante:
	cd template/project && make frizzante

configure:
	cd template/project && make configure

test:
	go run main.go -y --configure --app="template/project/app"
	test -d template/project/app && \
	rm -fr app && \
	mkdir -p app && \
	cp -r template/project/app .
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
	go run main.go --install
	cd template/project && make install

update:
	go run main.go --update
	cd template/project && make update

hooks:
	go run main.go --hooks

zip:
	./zip.sh

publish: zip
	./publish.sh