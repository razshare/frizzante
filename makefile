test: clean configure
	./test.sh

coverage: clean configure
	./coverage.sh

publish: test
	./publish.sh

configure:
	./configure.sh

clean:
	./clean.sh

install:
	./install.sh