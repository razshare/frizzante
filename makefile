test:
	CGO_ENABLED=1 go test

generate:
	go run lib/tools/aot/main.go
	bunx vite build --ssr .frz/router/server.ts --outDir .dist/server --emptyOutDir
	bunx vite build --outDir .dist/client --emptyOutDir

clean:
	go clean
	rm bin/app -fr
	rm node_modules -fr
	rm .dist -fr
	rm .frz -fr

update:
	go mod tidy
	bun update

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make clean && make update && make generate && make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit