build:
	go build -o build/siebes-gallery-downloader

build-m1:
	GOOS=darwin GOARCH=arm64 go build -o build/siebes-gallery-downloader-m1
