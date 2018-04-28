all: test web/dist
test:
	go test ./...
web/dist:
	cd web && npm run build
clean:
	rm -rf web/dist