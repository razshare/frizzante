test: clean setup
	./scripts/test.sh

coverage: clean setup
	./scripts/coverage.sh

publish: test
	./scripts/publish.sh

setup:
	./scripts/setup.sh

clean:
	./scripts/clean.sh

install: clean
	./scripts/install.sh
	make setup