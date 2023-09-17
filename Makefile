all: clean build run

push:
	git add .
	git commit -m "update"
	git push origin main

run:
	./mdocker

build:
	go build  -o mdocker main.go

clean:
	rm -rf mdocker

