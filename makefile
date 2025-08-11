frizzante:
	frizzante --app="template/app"

configure:
	go run main.go --configure --app="template/app"

test:
	go run main.go --test --app="template/app"

dev:
	go run main.go --dev --app="template/app"

package-watch:
	go run main.go --package-watch --app="template/app"

package:
	go run main.go --package --app="template/app"

check:
	go run main.go --check --app="template/app"

clean:
	go run main.go --clean --app="template/app"

format:
	go run main.go --format --app="template/app"

install:
	go run main.go --install --app="template/app"

update:
	go run main.go --update --app="template/app"

hooks:
	go run main.go --hooks --app="template/app"

publish:
	./publish.sh