build:
	go build -o morse led.go morse.go 
dependencies:
	go get github.com/pebbe/zmq4
	go get github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1
	go get github.com/golang/protobuf/proto