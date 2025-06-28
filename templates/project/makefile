########################
###### Composites ######
########################
test: configure-bun install check package
	CGO_ENABLED=1 go test ./...

build: configure-bun install check package
	CGO_ENABLED=1 go build -o .gen/bin/app .

dev: configure-air configure-bun install
	mkdir .gen/tmp -p
	mkdir app/dist -p
	touch app/dist/.gitkeep
	touch app/dist/server.js
	DEV=1 CGO_ENABLED=1 ./.gen/bin/air & \
	make package-watch & \
	wait

check: touch configure-bun
	cd app && \
	../.gen/bin/bun x eslint . && \
	../.gen/bin/bun x svelte-check --tsconfig ./tsconfig.json

package-watch: touch configure-bun
	cd app && \
	../.gen/bin/bun x vite build --logLevel info --ssr lib/utilities/frz/scripts/server.ts --outDir dist --watch & \
	cd app && \
	../.gen/bin/bun x vite build --logLevel info --outDir dist/client --watch & \
	wait

package: touch configure-bun
	cd app && \
	../.gen/bin/bun x vite build --logLevel info --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../.gen/bin/bun x vite build --logLevel info --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite && \
	touch dist/.gitkeep

install: touch configure-bun
	go mod tidy
	cd app && \
	../.gen/bin/bun install

update: touch configure-bun
	cd app && \
	../.gen/bin/bun update

format: touch configure-bun
	cd app && \
	../.gen/bin/bun x prettier --write .

########################
###### Primitives ######
########################
clean:
### Remove...
	go clean
	rm app/dist -fr
	rm app/node_modules -fr
	make touch

touch:
### Initialize...
	mkdir app/dist -p
	touch app/dist/.gitkeep
	touch app/dist/server.js

restore-frizzante-utilities:
	go run github.com/razshare/frizzante -u "app/lib/utilities"

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

configure-air:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make .gen/bin...
	mkdir .gen/bin -p
	# Get air...
	which .gen/bin/air || (curl -fsSL https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64 -o .gen/bin/air)
	chmod +x .gen/bin/air

configure: configure-bun configure-air
