all: test web/dist
test:
	go test ./...
web/dist:
	cd web && npm run build
dist/amd64: web/dist
	export GOOS=linux
	export GOARCH=amd64
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o koala cmd/koala/main.go
	tar czf dist/koala-linux-amd64.tar.gz koala web/dist
	rm koala
dist/arm: web/dist
	mkdir -p dist
	GOOS=linux GOARCH=arm go build -o koala cmd/koala/main.go
	tar czf dist/koala-linux-arm.tar.gz koala web/dist
	rm koala
dist: dist/amd64 dist/arm
clean:
	rm -rf web/dist dist/