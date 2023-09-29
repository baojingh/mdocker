all: clean build cp run

push:
	git add .
	git commit -m "update"
	git push origin main

cp:
	cp -r front/ /home/${USER}

run:
	./mdocker

build:
	go build  -o mdocker main.go

clean:
	rm -rf mdocker

