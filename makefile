########################
###### Composites ######
########################
test: install check package
	CGO_ENABLED=1 go test

publish: archive
### Publish...
	chmod +x ./publish.sh
	./publish.sh

########################
###### Primitives ######
########################
archive:
### Clean existing archives...
	rm project.zip -fr
	rm utilities.zip -fr
### Clean project template...
	rm templates/project/.gen -fr
	rm templates/project/.idea -fr
	rm templates/project/*.iml -fr
	cd templates/project && make clean
### Zip the project...
	cd templates/project && zip -9r ../../project.zip ./*
### Zip the utilities...
	cd templates/project/app/lib/utilities/frizzante && zip -9r ../../../../../../utilities.zip ./*

check:
	cd templates/project && make check

package:
	cd templates/project && make package

install:
	go mod tidy
	cd templates/project && make install

update:
	go mod tidy
	cd templates/project && make update

format:
	cd templates/project && make format

clean:
### Remove generated files...
	go clean
	rm .gen/out -fr
### Initialize template project...
	cd templates/project && make clean

hooks:
	printf "#!/usr/bin/env bash\n" > .git/hooks/pre-commit
	printf "make test" >> .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
