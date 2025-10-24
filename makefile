test: clean initialize
	./scripts/test.sh

coverage: clean initialize
	./scripts/coverage.sh

publish: test
	./scripts/publish.sh

initialize:
	./scripts/initialize.sh

clean:
	./scripts/clean.sh

install: clean
	./scripts/install.sh
	make initialize

package:
	./scripts/package.sh

dev:
	./scripts/dev.sh

types:
	./scripts/types.sh