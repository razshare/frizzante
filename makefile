build:
	make update
	make generate
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o bin/frizzante-cli-linux-386 .
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/frizzante-cli-linux-amd64 .
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -o bin/frizzante-cli-windows-386 .
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/frizzante-cli-windows-amd64 .

test:
	make update
	make generate
	CGO_ENABLED=1 go test ./frz

generate:
	make update
	go run main.go -generate -router -views="frz/components/views" -out=".frz/router"
	go run main.go -generate -utilities -out=".frz/utilities"
	bunx vite build --ssr .frz/router/server.ts --outDir frz/.dist/server --emptyOutDir
	bunx vite build --outDir frz/.dist/client --emptyOutDir
	rm .frz -fr

clean:
	go clean
	rm bin/app -fr
	rm node_modules -fr
	rm frz/.dist -fr
	rm .frz -fr

update:
	go mod tidy
	cd frz && go mod tidy
	bun update

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit