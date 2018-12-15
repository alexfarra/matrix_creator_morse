build:
	go build -o morse led.go morse.go 
dependencies:
	go get github.com/pebbe/zmq4
	go get github.com/matrix-io/matrix-protos-go