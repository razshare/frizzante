########################
###### Composites ######
########################
test: check package
	CGO_ENABLED=1 go test

check: configure-bun install
	cd app && \
	../.gen/bin/bun x eslint . && \
	../.gen/bin/bun x svelte-check --tsconfig ./tsconfig.json

package: configure-bun install
	cd app && \
	../.gen/bin/bun x vite build --logLevel info --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../.gen/bin/bun x vite build --logLevel info --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite && \
	touch dist/.gitkeep


install: configure-bun
	go mod tidy
	cd app && \
	../.gen/bin/bun install

update: configure-bun
	cd app && \
	../.gen/bin/bun update

format: configure-bun
	cd app && \
	../.gen/bin/bun x prettier --write .

########################
###### Primitives ######
########################
build:
	# Make .gen/out...
	mkdir .gen/out -p
	GOOS=linux GOARCH=amd64 go build -o ".gen/out/frizzante" && \
	cd .gen/out/ && \
	zip -9 frizzante-linux-am64.zip frizzante && \
	rm frizzante

clean:
### Remove...
	go clean
	rm .gen/out -fr
	rm app/dist -fr
	rm app/node_modules -fr
### Initialize...
	mkdir app/dist -p
	touch app/dist/.gitkeep
	touch app/dist/server.js

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

configure-bun:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make .gen/bin...
	mkdir .gen/bin -p
	# Get bun...
	which .gen/bin/bun || (curl -fsSL https://github.com/oven-sh/bun/releases/download/bun-v1.2.16/bun-linux-x64.zip -o .gen/bin/bun.zip && \
	unzip -j .gen/bin/bun.zip -d .gen/bin && rm .gen/bin/bun.zip -f)
	chmod +x .gen/bin/bun

configure: configure-bun
