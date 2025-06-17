###### Composites ######
build: clean configure-bun configure-air
	# Build for Linux Amd64...
	CGO_ENABLED=1 GOARCH="amd64" GOOS=linux go build -o bin/frizzante main.go && \
	zip bin/frizzante-amd64.zip bin/frizzante && \
	rm bin/frizzante -f

test: configure-bun generate check package
	CGO_ENABLED=1 cd frz && go test

update: configure-bun
	go mod tidy
	cd frz/app && \
	../../cli/bin/bun update

check: configure-bun
	cd frz/app && \
	../../cli/bin/bun x eslint . && \
	../../cli/bin/bun x svelte-check --tsconfig ./tsconfig.json

package: configure-bun
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	cd frz/app && \
	../../cli/bin/bun x vite build --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../../cli/bin/bun x vite build --outDir dist/client --emptyOutDir && \
	node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite

###### Primitives ######

generate:
	# Generate utilities...
	rm app/lib/utilities/frz -fr
	go run main.go -utilities -out="frz/app/lib/utilities/frz"

format:
	cd frz/app && \
	../../cli/bin/bun x prettier --write .

clean:
	go clean
	rm bin -fr
	rm cli/bin -fr
	rm frz/app/dist -fr
	rm frz/app/lib/utilities/frz -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	rm frz/app/node_modules -fr

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

configure-bun:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make cli/bin...
	mkdir cli/bin -p
	# Get bun...
	which cli/bin/bun || (curl -fsSL https://github.com/oven-sh/bun/releases/download/bun-v1.2.16/bun-linux-x64.zip -o cli/bin/bun.zip && \
	unzip -j cli/bin/bun.zip -d cli/bin && rm cli/bin/bun.zip -f)
	chmod +x cli/bin/bun

configure-air:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make cli/bin...
	mkdir cli/bin -p
	# Get air...
	which cli/bin/air || (curl -fsSL https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64 -o cli/bin/air)
	chmod +x cli/bin/air