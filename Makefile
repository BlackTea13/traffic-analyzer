.DEFAULT_GOAL := run

build: clean
	@ go build -o build/bin/app

run: build
	@ ./build/bin/app

clean:
	@ rm -rf ./build
