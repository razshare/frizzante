test: configure
	CGO_ENABLED=1 go test && \
	CGO_ENABLED=1 go test ./templates/pages

configure: clean update
	go run lib/prepare/main.go
	make www-build-server & \
	make www-build-client & \
	wait

clean:
	go clean
	rm bin/app -f
	rm cert.pem -f
	rm key.pem -f
	rm node_modules -fr
	rm .sessions -fr
	rm lib/config/.dist -fr
	rm .frizzante -fr
	mkdir lib/config/.dist/server -p
	mkdir lib/config/.dist/client -p
	touch lib/config/.dist/.gitkeep
	touch lib/config/.dist/server/.gitkeep
	touch lib/config/.dist/client/.gitkeep

update:
	go mod tidy
	bun update

www-build-server:
	bunx vite build --ssr .frizzante/vite-project/render.server.js --outDir lib/config/.dist/server --emptyOutDir && \
	./node_modules/.bin/esbuild lib/config/.dist/server/render.server.js --bundle --outfile=lib/config/.dist/server/render.server.js --format=esm --allow-overwrite

www-build-client:
	bunx vite build --outDir lib/config/.dist/client --emptyOutDir

certificate-interactive:
	openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem

certificate:
	openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem -nodes -subj \
	"/C=XX/ST=Test/L=Test/O=Test/OU=Test/CN=Test"

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

