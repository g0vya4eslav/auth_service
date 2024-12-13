clean:
	rm -f auth_service

build:
	go build -o auth_service main.go

run: build
	./auth_service

rebuild: clean build
