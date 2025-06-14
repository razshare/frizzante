test: configure generate update check package
	CGO_ENABLED=1 cd frz && go test

build: configure
	# Check requirements...
	command -v unzip >/dev/null || error 'zip is required to build frizzante'
	# Update...
	go mod tidy
	# Make bin...
	mkdir bin -p
	# Cleanup...
	rm bin/frizzante* -f
	# Build...
	CGO_ENABLED=1 GOARCH="amd64" GOOS=linux go build -o bin/frizzante main.go && \
	zip bin/frizzante-amd64.zip bin/frizzante && \
	rm bin/frizzante -f

update:
	go mod tidy
	cd frz/app && \
	../../bin/bun update

check:
	cd frz/app && \
	../../bin/bun x eslint . && \
	../../bin/bun x svelte-check --tsconfig ./tsconfig.json

package:
	rm frz/app/dist -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	cd frz/app && \
	../../bin/bun x vite build --ssr lib/utilities/frz/scripts/server.ts --outDir dist --emptyOutDir && \
	../../bin/bun x vite build --outDir dist/client --emptyOutDir

configure:
	# Check requirements...
	command -v unzip >/dev/null || error 'unzip is required to install and configure dependencies'
	command -v curl >/dev/null || error 'curl is required to install and configure dependencies'
	# Make bin...
	mkdir bin -p
	# Get bun...
	which bin/bun || \
	(curl -fsSL https://github.com/oven-sh/bun/releases/download/bun-v1.2.16/bun-linux-x64.zip -o bin/bun.zip && \
	unzip -j bin/bun.zip -d bin && rm bin/bun.zip -f)
	chmod +x bin/bun

generate: configure
	# Generate utilities...
	rm app/lib/utilities/frz -fr
	go run main.go -generate -utilities -out="frz/app/lib/utilities/frz"

format:
	cd frz/app && \
	../../bin/bun x prettier --write .

clean:
	go clean
	rm bin -fr
	rm frz/app/dist -fr
	rm frz/app/lib/utilities/frz -fr
	mkdir frz/app/dist/client -p
	touch frz/app/dist/client/index.html
	rm frz/app/node_modules -fr

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit