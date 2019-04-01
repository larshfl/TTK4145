package main

import (
	"flag"
	"fmt"
	"strconv"

	"./distributor"
	"./driver"
	"./network"
	statemachine "./stateMachine"
	"./types"

	//"os/exec"
	"runtime"
	"time"

	"./network/peers"
	"./setup"
	//"strconv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Example of how this works: go run main.go -port=10001"
	port := flag.String("port", "12345", "count of iterations")
	id := flag.String("id", "0", "id number")
	flag.Parse()
	portNum := *port
	ID := *id
	fmt.Printf("The port number is: %v\n", portNum)
	//(exec.Command("gnome-terminal", "-x", "sh", "-c", fmt.Sprintf("./SimElevatorServer --port %v",portNum))).Run()
	time.Sleep(2 * time.Second)

	// ID int
	intID, _ := strconv.Atoi(ID)

	setup.Init(portNum)

	var lightCh = make(chan []types.Elevator)
	go distributor.ChanUpdateLight(lightCh, intID)

	// Distributor <-> stateMachine
	var currentFloorCh = make(chan int, 1000)
	var directionCh = make(chan types.MotorDirection, 1000)
	var motorErrorCh = make(chan bool, 1000)
	var completedOrdersCh = make(chan types.SingleOrder, 1000)
	var orderListCh = make(chan []types.SingleOrder)

	// StateMachine <-> driver
	var floorArrivalsCh = make(chan int)

	// Distributor <-> driver
	var buttonEventCh = make(chan types.ButtonEvent, 1000)

	//Distributor <-> Network
	var ElevToNetCh = make(chan []types.Elevator, 16)
	var ElevToDistrCh = make(chan []types.Elevator, 16)
	var networkEnableCh = make(chan bool)
	var singleOrderCh = make(chan types.SingleOrder)
	var elevOnNetworkCh = make(chan peers.PeerUpdate)
	var myIDCh = make(chan string)

	go statemachine.StateMachine(currentFloorCh,
		directionCh,
		motorErrorCh,
		completedOrdersCh,
		orderListCh,
		floorArrivalsCh)

	go distributor.Distributor(currentFloorCh,
		buttonEventCh,
		elevOnNetworkCh,
		completedOrdersCh,
		directionCh,
		motorErrorCh,
		ElevToNetCh,
		networkEnableCh,
		orderListCh,
		singleOrderCh,
		ElevToDistrCh,
		myIDCh,
		lightCh)

	go network.Network(ElevToNetCh,
		ElevToDistrCh,
		networkEnableCh,
		elevOnNetworkCh,
		myIDCh,
		portNum,
		ID)

	go driver.PollButtons(buttonEventCh)
	go driver.PollFloorSensor(floorArrivalsCh)

	//test section

	// type Elev struct{
	// 	Name string
	// 	Age int
	// }

	// 	e := types.Elevator{
	// 		Floor: 0,
	// 		Dir: types.MotorDirectionStop,
	// 		Orders: [4][3]int{{0, 0, 0},{1, 0 ,0},{0, 0, 0}, {0, 0, 0}},
	// 		Behaviour: types.Idle,
	// 		ID: "",
	// 	}
	// 	fmt.Print(e)
	// 	fmt.Print("\n")

	// 	m := make(map[int]*types.Elevator)
	// 	m[1] = &e
	// 	m[1].Floor = 3

	// 	fmt.Print(m[1])

	select {}
}
