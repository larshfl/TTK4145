//driver module

package driver

import "time"
import "sync"
import "net"
import "fmt"
import "../types"



const _pollRate = 20 * time.Millisecond

var _initialized = false
var _numFloors = 4
var _mtx sync.Mutex
var _conn net.Conn

// Public funtions

// Init initializes the driver and opens an connection through TCP
func Init(addr string, numFloors int) {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	_initialized = true
}

// SetMotorDirection sets the motor direction - Up, Down and Stop
func SetMotorDirection(dir types.MotorDirection) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{1, byte(dir), 0, 0})
}

// SetButtonLamp sets the lamps - Hall up, Hall down and Cab
func SetButtonLamp(button types.ButtonType, floor int, value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{2, byte(button), byte(floor), toByte(value)})
}

// SetFloorIndicator sets the Floor indicator
func SetFloorIndicator(floor int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{3, byte(floor), 0, 0})
}

// SetDoorOpenLamp - to turn on or off the door open lamp
func SetDoorOpenLamp(value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{4, toByte(value), 0, 0})
}

// SetStopLamp - to turn on or off the Stop lamp
func SetStopLamp(value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{5, toByte(value), 0, 0})
}

// PollButtons polls ButtonType - HallUP, HallDown and Cab
func PollButtons(receiver chan<- types.ButtonEvent) { 
	prev := make([][3]bool, _numFloors)
	for {
		time.Sleep(_pollRate)
		for f := 0; f < _numFloors; f++ {
			for b := types.ButtonType(0); b < 3; b++ {
				v := getButton(b, f)
				if v != prev[f][b] && v != false {
					receiver <- types.ButtonEvent{Floor: f, Button: b}
				}
				prev[f][b] = v
			}
		}
	}
}

// PollFloorSensor polls the floor sensor - Send an int for current floor/floor pased
func PollFloorSensor(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(_pollRate)
		v := getFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}



// Internal functions

func getButton(button types.ButtonType, floor int) bool {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{6, byte(button), byte(floor), 0})
	var buf [4]byte
	_conn.Read(buf[:]) 
	return toBool(buf[1])
}

func getFloor() int {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{7, 0, 0, 0})
	var buf [4]byte
	_conn.Read(buf[:])
	if buf[1] != 0 {
		return int(buf[2])
	} else {
		return -1
	}
}

func toByte(a bool) byte {
	var b byte
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b = false
 	if a != 0 {
		b = true
	}
	return b
}


