all: clean build run

push:
	git add .
	git commit -m "update"
	git push origin main

run:
	./mdocker

build:
	go build  -o mdocker main.go

term:
	sudo kill -TERM  $(ps -ef | grep mdocker | grep -v "grep" | awk '{print $2}')

kill:
	sudo kill -9     $(ps -ef | grep mdocker | grep -v "grep" | awk '{print $2}')


clean:
	rm -rf mdocker

