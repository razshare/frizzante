########################
###### Composites ######
########################
test: check package
	CGO_ENABLED=1 go test

check: configure-bun install
	cd app && \
	../.gen/bin/bun x eslint . && \
	../.gen/bin/bun x svelte-check --tsconfig ./tsconfig.json
	cd templates/project && make check

package: configure-bun install
	cd app && \
	../.gen/bin/bun x vite build --logLevel info --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../.gen/bin/bun x vite build --logLevel info --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite && \
	touch dist/.gitkeep

install: configure-bun
	go mod tidy
	cd app && ../.gen/bin/bun install
	cd templates/project && make install

update: configure-bun
	cd app && ../.gen/bin/bun update
	cd templates/project && make update

format: configure-bun
	cd app && ../.gen/bin/bun x prettier --write .
	cd templates/project && make format

########################
###### Primitives ######
########################
publish:
### Clean project template...
	rm templates/project/.gen -fr
	rm templates/project/.idea -fr
	rm templates/project/*.iml -fr
	cd templates/project && make clean
### Zip the project...
	cd templates/project && zip -9r ../../project.zip ./*
	chmod +x ./publish.sh
	./publish.sh

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
	cd templates/project && make clean

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
