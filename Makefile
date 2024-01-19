all: clean build run

run:
	./mdocker

build:
	go build  -o mdocker main.go

clean:
	rm -rf mdocker

