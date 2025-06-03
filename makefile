test:
	make generate
	CGO_ENABLED=1 go test ./frz

generate:
	frizzante-cli -generate -router -views="lib/components/views" -out=".frz/router"
	frizzante-cli -generate -utilities -out=".frz/utilities"
	bunx vite build --ssr .frz/router/server.ts --outDir frz/.dist/server --emptyOutDir
	bunx vite build --outDir frz/.dist/client --emptyOutDir
	rm .frz -fr

clean:
	go clean
	rm node_modules -fr
	rm cli/bin -fr
	rm cli/.frz -fr
	rm frz/.dist -fr

update:

	go mod tidy
	bun update

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit