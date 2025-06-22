###### Composites ######
test: configure-bun check package
	CGO_ENABLED=1 cd frz && go test

update: configure-bun
	go mod tidy
	cd frz/app && \
	../../bin/bun update

check: configure-bun
	cd frz/app && \
	../../bin/bun x eslint . && \
	../../bin/bun x svelte-check --tsconfig ./tsconfig.json

package: configure-bun
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	cd frz/app && \
	../../bin/bun x vite build --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../../bin/bun x vite build --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite

###### Primitives ######
format:
	cd frz/app && \
	../../bin/bun x prettier --write .

clean:
	go clean
	rm cli/starter.zip -fr
	rm bin -fr
	mkdir bin -p
	touch cli/.gitkeep
	rm frz/app/dist -fr
	rm frz/app/lib/utilities/frz -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	rm frz/app/node_modules -fr

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

configure-starter:
	which cli/starter.zip || \
	curl -fsSL https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip -o cli/starter.zip

configure-bun:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make bin...
	mkdir bin -p
	# Get bun...
	which bin/bun || (curl -fsSL https://github.com/oven-sh/bun/releases/download/bun-v1.2.16/bun-linux-x64.zip -o bin/bun.zip && \
	unzip -j bin/bun.zip -d bin && rm bin/bun.zip -f)
	chmod +x bin/bun

configure-air:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make bin...
	mkdir bin -p
	# Get air...
	which bin/air || (curl -fsSL https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64 -o bin/air)
	chmod +x bin/air

configure: clean configure-bun configure-air