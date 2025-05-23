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
	rm .dist -fr
	rm .frizzante -fr
	mkdir .dist/server -p
	mkdir .dist/client -p
	touch .dist/.gitkeep
	touch .dist/server/.gitkeep
	touch .dist/client/.gitkeep

update:
	go mod tidy
	bun update

www-build-server:
	bunx vite build --ssr .frizzante/vite-project/render.server.js --outDir .dist/server --emptyOutDir && \
	./node_modules/.bin/esbuild .dist/server/render.server.js --bundle --outfile=.dist/server/render.server.js --format=esm --allow-overwrite

www-build-client:
	bunx vite build --outDir .dist/client --emptyOutDir

certificate-interactive:
	openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem

certificate:
	openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem -nodes -subj \
	"/C=XX/ST=Test/L=Test/O=Test/OU=Test/CN=Test"

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

