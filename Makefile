all:
	./compile.rb

darwin:
	go build -o bin/star_darwin-amd64
	go build -o bin/star_darwin-386

linux:
	go build -o bin/star_linux-amd64
	go build -o bin/star_linux-386

windows:
	go build -o bin/star_windows-amd64
	go build -o bin/star_windows-386
