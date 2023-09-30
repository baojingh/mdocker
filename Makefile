all: clean build cp run

push:
	git add .
	git commit -m "update"
	git push origin main

cp:
	sudo rm -rf /var/www/html/front
	sudo cp -r front/ /var/www/html
	sudo systemctl  restart nginx

run:
	./mdocker

build:
	go build  -o mdocker main.go

clean:
	rm -rf mdocker

