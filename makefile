frizzante:
	frizzante

configure:
	go run main.go --configure

test:
	go run main.go -y --configure --app="template/project/app"
	test -d template/project/app && \
	rm -fr app && \
	mkdir -p app && \
	cp -r template/project/app .
	go run main.go -y --test
	rm -fr app


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

zip:
	./zip.sh

publish: zip
	./publish.sh