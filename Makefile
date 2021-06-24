all:
	go build -o dondu main.go

install:
	cp ./dondu /usr/local/bin

clean:
	rm ./dondu
