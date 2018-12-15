package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	if len(os.Args) < 2 {
		panic("No IP address given!")
	}
	wait.Add(1)
	go initPort()
	go dataUpdatePort()
	go errorPort()
	wait.Wait()
}

var wait sync.WaitGroup
var everloop = core.EverloopImage{}
var blankEverloop = core.EverloopImage{}
var configPusher *zmq.Socket
var config core.LedValue
var encodedString string
var sleepSlice []delay
var morseDelay = map[rune]int{
	'.': dotDelay,
	'-': dashDelay,
}

type delay struct {
	alive int
	dead  int
}

const (
	charOffDelay = 150
	wordOffDelay = charOffDelay * 3
	dotDelay     = 150
	dashDelay    = dotDelay * 3
)

// parses IP into RGBW values and morse code
func handleInput(ip string) {
	encodedString = encode(ip)
	rgbw := make([]uint32, 4)
	vals := strings.Split(ip, ".")
	if len(vals) != 4 {
		panic("Not a valid IP")
	}
	for i, val := range vals {
		if temp, err := strconv.Atoi(val); err == nil && temp < 256 {
			rgbw[i] = uint32(temp)
		} else {
			panic("Not a valid IP")
		}
	}
	config = core.LedValue{
		Red:   rgbw[0],
		Green: rgbw[1],
		Blue:  rgbw[2],
		White: rgbw[3],
	}
	everloop.Led = []*core.LedValue{}
	for i := int32(0); i < everloop.EverloopLength; i++ {
		everloop.Led = append(everloop.Led, &config)
	}
	initBlankLights()
}

// fills in the slice corresponding to the morse code delay sleep time
func generateDelays() {
	sleepSlice = make([]delay, 0, len(encodedString))
	previousSpace := false
	for i, rune := range encodedString {
		temp := delay{}
		if i == 0 {
			temp.alive = morseDelay[rune]
		} else {
			if rune == ' ' {
				sleepSlice[len(sleepSlice)-1].dead = wordOffDelay
				previousSpace = true
				continue
			} else {
				if !previousSpace {
					sleepSlice[len(sleepSlice)-1].dead = charOffDelay
				}
				previousSpace = false
				temp.alive = morseDelay[rune]
			}
		}
		sleepSlice = append(sleepSlice, temp)
	}
	if len(sleepSlice) > 0 {
		sleepSlice[len(sleepSlice)-1].dead = -1
	}
}

// Connect to the base port to send configurations
func basePort() {
	defer wait.Done()
	configPusher, _ = zmq.NewSocket(zmq.PUSH)
	configPusher.Connect("tcp://127.0.0.1:20021")
	generateDelays()
	for i := 0; i < len(sleepSlice); i++ {
		sendState(true)
		time.Sleep(time.Duration(sleepSlice[i].alive) * time.Millisecond)
		resetLights()
		if sleepSlice[i].dead != -1 {
			time.Sleep(time.Duration(sleepSlice[i].dead) * time.Millisecond)
		}
	}
}

// Sends current config from everloop.Led to device via given socket
func sendState(on bool) {
	var state core.DriverConfig
	if on {
		state = core.DriverConfig{
			Image: &everloop,
		}
	} else {
		state = core.DriverConfig{
			Image: &blankEverloop,
		}
	}

	encodedState, _ := proto.Marshal(&state)
	configPusher.Send(string(encodedState), 1)
}

// reset LEDs to off state
func resetLights() {
	sendState(false)
}

// Initialize the blank everloop config
func initBlankLights() {
	blankEverloop.Led = []*core.LedValue{}
	blankConfig := core.LedValue{
		Red:   0,
		Green: 0,
		Blue:  0,
		White: 0,
	}
	for i := int32(0); i < blankEverloop.EverloopLength; i++ {
		blankEverloop.Led = append(blankEverloop.Led, &blankConfig)
	}
}

// Create initial connection with Creator
func initPort() {
	pusher, _ := zmq.NewSocket(zmq.PUSH)
	pusher.Connect("tcp://127.0.0.1:20022")
	pusher.Send("", 1)
}

// Connect to error port
func errorPort() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect("tcp://127.0.0.1:20023")
	subscriber.SetSubscribe("")
	for {
		message, _ := subscriber.Recv(2)
		fmt.Println("ERROR:", message)
	}
}

// Connect to data update port, responsible for filling in the everloop struct
func dataUpdatePort() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect("tcp://127.0.0.1:20024")
	subscriber.SetSubscribe("")
	for {
		message, _ := subscriber.Recv(2)
		proto.Unmarshal([]byte(message), &everloop)
		blankEverloop = everloop
		handleInput(os.Args[1])
		go basePort()
		return
	}
}
