package main

import (
	"flag"
	"fmt"
	"strconv"

	"./distributor"
	"./driver"
	"./network"
	"./statemachine"
	"./types"

	"os/exec"
	"runtime"
	"time"

	"./network/peers"
	"./setup"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.String("port", "12345", "count of iterations")
	id := flag.String("id", "0", "id number")
	flag.Parse()
	portNum := *port
	myID := *id
	fmt.Printf("The port number is: %v\n", portNum)
	fmt.Printf("The id is: %v\n", myID)
	(exec.Command("gnome-terminal", "-x", "sh", "-c", fmt.Sprintf("./SimElevatorServer --port %v",portNum))).Run()
	time.Sleep(2 * time.Second)

	ID, _ := strconv.Atoi(myID)

	setup.Init(portNum)

	var lightCh = make(chan []types.Elevator)
	go driver.ChanUpdateButtonLights(lightCh, ID)

	// distributor <-> stateMachine
	var currentFloorCh 			= make(chan int,4)
	var directionCh 			= make(chan types.MotorDirection,3)
	var motorErrorCh		 	= make(chan bool)
	var completedOrdersCh 		= make(chan types.SingleOrder,3)
	var orderListCh 			= make(chan []types.SingleOrder)

	// dtateMachine <-> driver
	var floorArrivalsCh 		= make(chan int)

	// distributor <-> driver
	var buttonEventCh 			= make(chan types.ButtonEvent)

	// distributor <-> network
	// ElevToNetCh and ElevToDistrCh are buffered to 16 because this is the maximum number 
	// of orders that are possible to have in the FIFO queue simultaniously 
	var ElevToNetCh 			= make(chan []types.Elevator, 16)
	var ElevToDistrCh 			= make(chan []types.Elevator, 16)
	var networkEnableCh 		= make(chan bool)
	var singleOrderCh 			= make(chan types.SingleOrder)
	var elevOnNetworkCh 		= make(chan peers.PeerUpdate)

	go statemachine.StateMachine(
		currentFloorCh,
		directionCh,
		motorErrorCh,
		completedOrdersCh,
		orderListCh,
		floorArrivalsCh)

	go distributor.Distributor(
		currentFloorCh,
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
		lightCh, 
		ID)

	go network.Network(
		ElevToNetCh,
		ElevToDistrCh,
		networkEnableCh,
		elevOnNetworkCh,
		portNum,
		myID)

	go driver.PollButtons(buttonEventCh)
	go driver.PollFloorSensor(floorArrivalsCh)

	select {}
}
