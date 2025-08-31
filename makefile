zip:
	rm -fr internal/project/.gen
	rm -fr internal/project/app/dist
	rm -fr internal/project/app/.vite
	rm -fr internal/project/app/node_modules
	rm -fr internal/project/app/lib/core/svelte/ssr/app
	cd internal && zip -rq9 project.zip project

publish:
	./publish.sh