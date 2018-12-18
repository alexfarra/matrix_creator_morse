all: build

build:
	go build -o morse led.go morse.go
