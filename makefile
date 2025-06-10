test:
	make update
	make check
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	make package
	CGO_ENABLED=1 cd frz && go test

format:
	cd frz/app && \
	bunx prettier --write .

clean:
	go clean
	rm bin/app -fr
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	rm frz/app/node_modules -fr

update:
	go mod tidy
	cd frz/app && \
	bun update

check:
	cd frz/app && \
	bunx eslint . && \
	bunx svelte-check --tsconfig ./tsconfig.json

generate:
	rm frz/app/lib/utilities -fr
	go run main.go -generate -utilities -out="frz/app/lib/utilities"

package:
	cd frz/app && \
	bunx vite build --ssr lib/utilities/scripts/server.ts --outDir dist --emptyOutDir && \
	bunx vite build --outDir dist/client --emptyOutDir

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit