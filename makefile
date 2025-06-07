test:
	make update
	make generate
	rm frz/dist -fr
	make package
	CGO_ENABLED=1 cd frz && go test

clean:
	go clean
	rm bin/app -fr
	rm frz/dist -fr
	rm app/lib/utilities -fr
	rm node_modules -fr

update:
	go mod tidy
	bun update

generate:
	go run main.go -generate -utilities -out="app/lib/utilities"

package:
	bunx vite build --ssr app/lib/utilities/scripts/server.ts --outDir frz/dist/server --emptyOutDir
	bunx vite build --outDir frz/dist/client --emptyOutDir

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit