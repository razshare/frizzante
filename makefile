frizzante:
	frizzante

configure:
	go run main.go --configure

test:
	go run main.go --test

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

publish:
	./publish.sh