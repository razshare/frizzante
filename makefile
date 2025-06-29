########################
###### Composites ######
########################
test: install check package fadeout
	CGO_ENABLED=1 go test || make fadein
	make fadein

publish: fadein archive
### Publish...
	chmod +x ./publish.sh
	./publish.sh

########################
###### Primitives ######
########################
fadein:
	test -f templates/project/main.go || mv templates/project/main.go.txt templates/project/main.go
	test -f templates/project/go.mod || mv templates/project/go.mod.txt templates/project/go.mod
	test -f templates/project/go.sum || mv templates/project/go.sum.txt templates/project/go.sum

fadeout:
	test -f templates/project/main.go.txt || mv templates/project/main.go templates/project/main.go.txt
	test -f templates/project/go.mod.txt || mv templates/project/go.mod templates/project/go.mod.txt
	test -f templates/project/go.sum.txt || mv templates/project/go.sum templates/project/go.sum.txt

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
	cd templates/project && zip -9r ../../project.zip * .[^.]*
### Zip the utilities...
	cd templates/project/app/lib/frizzante && zip -9r ../../../../../../utilities.zip * .[^.]*

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
