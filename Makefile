remove-bind-go:
	rm -f bind.go

bundle: remove-bind-go
	astilectron-bundler -v
	rm bind_*.go
	cp .bind.go.src bind.go
