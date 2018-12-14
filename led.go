package main

import (
	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	go initPort()
	go dataUpdatePort()
}

var everloop = core.EverloopImage{}
var configPusher *zmq.Socket

// Connect to the base port to send configurations
func basePort() {
	configPusher, _ = zmq.NewSocket(zmq.PUSH)
	configPusher.Connect("tcp://127.0.0.1:20021")
	resetLights()
}

// Sends current config from everloop.Led to device via given socket
func sendState() {
	state := core.DriverConfig{
		Image: &everloop,
	}

	encodedState, _ := proto.Marshal(&state)
	configPusher.Send(string(encodedState), 1)
}

// reset LEDs to off state
func resetLights() {
	everloop.Led = []*core.LedValue{}
	config := core.LedValue{
		Red:   0,
		Green: 0,
		Blue:  0,
		White: 0,
	}
	for i := int32(0); i < everloop.EverloopLength; i++ {
		everloop.Led = append(everloop.Led, &config)
	}

	sendState()
}

func initPort() {
	pusher, _ := zmq.NewSocket(zmq.PUSH)
	pusher.Connect("tcp://127.0.0.1:20022")
	pusher.Send("", 1)
}

func dataUpdatePort() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect("tcp://127.0.0.1:20024")
	subscriber.SetSubscribe("")
	for {
		message, _ := subscriber.Recv(2)
		proto.Unmarshal([]byte(message), &everloop)
		go basePort()
		return
	}
}
