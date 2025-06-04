test:
	make update
	make generate
	rm frz/dist -fr
	bunx vite build --ssr app/frz/scripts/server.ts --outDir frz/dist/server --emptyOutDir
	bunx vite build --outDir frz/dist/client --emptyOutDir
	CGO_ENABLED=1 cd frz && go test

generate:
	go run main.go -generate -utilities -out="app/frz"

update:
	go mod tidy
	bun update

clean:
	go clean
	rm bin/app -fr
	rm frz/dist -fr
	rm app/frz -fr
	rm node_modules -fr

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit