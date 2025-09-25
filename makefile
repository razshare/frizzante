test: clean configure package
	./scripts/test.sh

coverage: clean configure package
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

package:
	./scripts/package.sh

dev:
	./scripts/dev.sh

types:
	./scripts/types.sh