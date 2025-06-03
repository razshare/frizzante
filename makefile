build:
	make update
	GOOS=linux GOARCH=386 cd cli && go build -o bin/frizzante-cli-linux-386 .
	GOOS=linux GOARCH=amd64 cd cli && go build -o bin/frizzante-cli-linux-amd64 .
	GOOS=windows GOARCH=386 cd cli && go build -o bin/frizzante-cli-windows-386 .
	GOOS=windows GOARCH=amd64 cd cli && go build -o bin/frizzante-cli-windows-amd64 .

test:
	make generate
	CGO_ENABLED=1 cd frz && go test

generate:
	make update
	cd cli && go run main.go -generate -router -views="lib/components/views" -out=".frz/router"
	cd cli && go run main.go -generate -utilities -out=".frz/utilities"
	bunx vite build --ssr cli/.frz/router/server.ts --outDir frz/.dist/server --emptyOutDir
	bunx vite build --outDir frz/.dist/client --emptyOutDir
	rm .frz -fr

clean:
	cd cli && go clean
	cd frz && go clean
	rm node_modules -fr
	rm cli/bin -fr
	rm cli/.frz -fr
	rm frz/.dist -fr

update:
	cd cli && go mod tidy
	cd frz && go mod tidy
	bun update

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit