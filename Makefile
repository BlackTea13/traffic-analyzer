.DEFAULT_GOAL := run

build:
	@ go build -o build/bin/app

run: build
	@ ./build/bin/app
