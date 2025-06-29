########################
###### Composites ######
########################
test: install check package
	CGO_ENABLED=1 go test ./...

build: install check package
	CGO_ENABLED=1 go build -o .gen/bin/app .

dev: configure-air install
	mkdir .gen/tmp -p
	mkdir app/dist -p
	touch app/dist/.gitkeep
	touch app/dist/server.js
	DEV=1 CGO_ENABLED=1 ./.gen/bin/air & \
	make package-watch & \
	wait

check: touch
	cd app && \
	bunx eslint . && \
	bunx svelte-check --tsconfig ./tsconfig.json

package-watch: touch
	cd app && \
	bunx vite build --logLevel info --ssr lib/utilities/frizzante/scripts/server.ts --outDir dist --watch & \
	cd app && \
	bunx vite build --logLevel info --outDir dist/client --watch & \
	wait

package: touch
	cd app && \
	bunx vite build --logLevel info --ssr lib/utilities/frizzante/scripts/server.ts --outDir dist --emptyOutDir && \
	bunx vite build --logLevel info --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite && \
	touch dist/.gitkeep

install: touch
	go mod tidy
	cd app && \
	bun install

update: touch
	cd app && \
	bun update

format: touch
	cd app && \
	bunx prettier --write .

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
	go run github.com/razshare/frizzante -r "app/lib/utilities/frizzante"

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

configure-air:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make .gen/bin...
	mkdir .gen/bin -p
	# Get air...
	which .gen/bin/air || (curl -fsSL https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64 -o .gen/bin/air)
	chmod +x .gen/bin/air

configure:
