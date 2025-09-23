test: clean configure
	./scripts/test.sh

coverage: clean configure
	./scripts/coverage.sh

publish: test
	./scripts/publish.sh

configure:
	./scripts/configure.sh

clean:
	./scripts/clean.sh

install: clean
	./scripts/install.sh
	make configure

dev:
	./scripts/dev.sh

types:
	./scripts/types.sh