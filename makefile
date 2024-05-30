build:
	@GOOS=windows GOARCH=amd64 go build -o bin/windows/main.exe cmd/server.go

build_linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/main cmd/server.go

run: build
	@./bin/windows/main.exe -port 8080 -dbengine edgedb -edgedb_instance Ghostship

run_linux: build_linux
	@./bin/linux/main
