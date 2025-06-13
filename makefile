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
	../../bin/bun x prettier --write .

clean:
	go clean
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	rm frz/app/node_modules -fr

update:
	go mod tidy
	cd frz/app && \
	../../bin/bun update

check:
	cd frz/app && \
	../../bin/bun x eslint . && \
	../../bin/bun x svelte-check --tsconfig ./tsconfig.json

package:
	cd frz/app && \
	../../bin/bun x vite build --ssr lib/utilities/scripts/server.ts --outDir dist --emptyOutDir && \
	../../bin/bun x vite build --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite

configure:
	mkdir bin -p
	which bin/bun || \
	(curl -fsSL https://github.com/oven-sh/bun/releases/latest/download/bun-linux-x64.zip -o bin/bun.zip && \
	unzip -j bin/bun.zip -d bin && rm bin/bun.zip -f)
	which bin/air || curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s
	go run main.go -generate -utilities -out="frz/app/lib/utilities"


hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit