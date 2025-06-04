test:
	CGO_ENABLED=1 go test ./...

generate:
	go run main.go -generate -render -views="frz/lib/components/views" -out="frz/.generated/render"
	go run main.go -generate -utilities -out="frz/.generated/utilities"
	bunx vite build --ssr frz/.generated/render/server.ts --outDir frz/.dist/server --emptyOutDir
	bunx vite build --outDir frz/.dist/client --emptyOutDir

clean:
	go clean
	rm node_modules -fr
	rm frz/.generated -fr
	rm frz/.dist -fr

update:
	go mod tidy
	bun update

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit