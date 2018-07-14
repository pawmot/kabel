remove-bind-stub:
	rm -f bind.go

restore-bind-stub:
	rm bind*.go
	cp .bind.go.src bind.go

bundle: remove-bind-stub
	astilectron-bundler -v
	$(MAKE) restore-bind-stub
