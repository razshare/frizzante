test:
	./test.sh || make lock

coverage:
	./coverage.sh

lock:
	(test -f main.txt && mv main.txt main.go) &
	mv internal/project/go.mod.txt internal/project/go.mod &
	mv internal/project/go.sum.txt internal/project/go.sum &
	wait

unlock:
	(test -f main.txt && mv main.txt main.go) &
	mv internal/project/go.mod internal/project/go.mod.txt &
	mv internal/project/go.sum internal/project/go.sum.txt &
	wait

publish: test
	./publish.sh || make lock