push:
	git add .
	git commit -m "update"
	git push origin main

build:
	go build  -o mdocker main.go

install:
	cp mdocker /usr/local/bin

clean:
	rm -rf mdocker



