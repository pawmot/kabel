.DEFAULT_GOAL := bundle-dev

remove-bind-stub:
	rm -f bind.go

restore-bind-stub:
	rm -f bind*.go
	cp .bind.go.src bind.go

frontend-dist:
	cd frontend && ng build --prod

frontend-dist-watch:
	cd frontend && ng build --watch

copy-frontend-dist:
	rm -f resources/app/*
	cp frontend/dist/kabel/* resources/app

bundle: frontend-dist copy-frontend-dist remove-bind-stub
	astilectron-bundler -v
	$(MAKE) restore-bind-stub

# This target assumes that `ng build` has been run externally.
# The indended way is to run `ng build --watch` that greatly speeds up the process in comparison to a full build.
# You can use `frontend-dist-watch` target to run the aformentioned build.
bundle-dev: remove-bind-stub copy-frontend-dist
	astilectron-bundler -v
	$(MAKE) restore-bind-stub