push:
	git add .
	git commit -m "update"
	git push origin main


build:
	docker build .

install:
	cp ./ /usr/bin

clean:
	rm -rf ./rr



