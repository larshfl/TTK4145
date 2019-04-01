//State Machine

package statemachine

import (
	"../driver"
	"../types"
	"./requests"
)

// StateMachine for a single Elevator
func StateMachine(
	currentFloorCh 		chan int, 
	directionCh 		chan types.MotorDirection,
	motorErrorCh 		chan bool, 
	completedOrdersCh 	chan types.SingleOrder,
	orderListCh 		chan []types.SingleOrder, 
	floorArrivalsCh 	chan int) {

	var e types.Elevator

	doorTimerFinished := make(chan bool)
	newDoorTimer := make(chan bool)
	go requests.OpenDoorTimer(doorTimerFinished, newDoorTimer)

	motorError := make(chan bool)
	resetMotorTimer := make(chan bool)
	go requests.CheckForMotorError(motorError, resetMotorTimer, &e)

	//Drive the elevator to the initial position - Floor 0
	driver.SetMotorDirection(types.MotorDirectionDown)
	for {
		e.Floor = <-floorArrivalsCh
		driver.SetFloorIndicator(e.Floor)
		currentFloorCh <- e.Floor
		if e.Floor > 0 {
			driver.SetMotorDirection(types.MotorDirectionDown)
		} else {
			driver.SetMotorDirection(types.MotorDirectionStop)
			e.Behaviour = types.Idle
			break
		}
	}

	for {
		select {
		case requests.OrderList = <-orderListCh:
			switch e.Behaviour {
			case types.Idle:
				if e.Floor == requests.OrderList[0].Floor {
					e.Behaviour = types.DoorOpen
					newDoorTimer <- true
					driver.SetDoorOpenLamp(true)
				} else {
					e.Dir = requests.ChooseDirection(e)
					directionCh <- e.Dir
					driver.SetMotorDirection(e.Dir)
					e.Behaviour = types.Moving
					resetMotorTimer <- true
				}
			case types.Moving:
			case types.DoorOpen:
			case types.Undefined:
			}

		case e.Floor = <-floorArrivalsCh:
			driver.SetFloorIndicator(e.Floor)
			currentFloorCh <- e.Floor
			resetMotorTimer <- true
			switch e.Behaviour {
			case types.Idle:
			case types.Moving:
				if requests.ShouldStop(e) {
					driver.SetMotorDirection(types.MotorDirectionStop)
					newDoorTimer <- true
					e.Behaviour = types.DoorOpen
					driver.SetDoorOpenLamp(true)
				}
			case types.DoorOpen:
			case types.Undefined:
			}

		case <-doorTimerFinished:
			switch e.Behaviour {
			case types.Idle:
			case types.Moving:
			case types.DoorOpen:
				requests.ClearOrders(completedOrdersCh, e)
				driver.SetDoorOpenLamp(false)
				e.Dir = requests.ChooseDirection(e)
				directionCh <- e.Dir
				driver.SetMotorDirection(e.Dir)
				if e.Dir == types.MotorDirectionStop {
					e.Behaviour = types.Idle
				} else {
					e.Behaviour = types.Moving
					resetMotorTimer <- true
				}
			case types.Undefined:
			}

		case <-motorError:
			motorErrorCh <- true
		}
	}
}
