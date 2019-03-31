//Package distributor is the module for fiksing the orders chaos
package distributor

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"../driver"
	"../network/peers"
	"../types"
)

var MotorError = false

// Distributor does..
func Distributor(currentFloorCh chan int,
	buttonEventCh chan types.ButtonEvent,
	elevOnNetworkCh chan peers.PeerUpdate,
	completedOrderCh chan types.SingleOrder,
	directionCh chan types.MotorDirection,
	motorErrorCh chan bool,
	ElevToNetCh chan []types.Elevator,
	orderConfirmCh chan int, turnOfNetworkCh chan bool,
	orderListCh chan []types.SingleOrder,
	singleOrderCh chan types.SingleOrder,
	ElevToDistrCh chan []types.Elevator,
	myIDCh chan string,
	lightCh chan []types.Elevator) {

	ElevSlice := make([]types.Elevator, 0)
	myIP := <-myIDCh
	StateMachineOrderSlice := make([]types.SingleOrder, 0)
	orderCount := 0
	var elevOnNet peers.PeerUpdate

	for {
		select {
		case floor := <-currentFloorCh:
			MotorError = false
			place := sToi(myIP, ElevSlice)
			ElevSlice[place].Floor = floor
			fmt.Printf("Received foor arrival: %v inElevSlice = %v \n", floor, ElevSlice[sToi(myIP, ElevSlice)].Floor)

		case buttonPress := <-buttonEventCh:
			fmt.Printf("in button Press \n")
			if ElevatorMotorError(buttonPress, motorErrorCh, turnOfNetworkCh) {
				break
			}
			if isDuplicate(buttonPress, ElevSlice, elevOnNet, myIP) {
				break
			}

			lowestCost := myIP
			min := math.Inf(1)
			for ipIndex := 0; ipIndex < len(elevOnNet.Peers); ipIndex++ {
				fmt.Printf("Elev nr: %v floor: %v \n", ElevSlice[ipIndex].ID, ElevSlice[ipIndex].Floor)

				cost := TimeToIdle(ElevSlice[ipIndex], buttonPress)
				if min > cost {
					min = cost
					lowestCost = elevOnNet.Peers[ipIndex]
				}
			}
			fmt.Printf("lowest cost: %v \n", lowestCost)

			if types.ButtonCab == buttonPress.Button {
				lowestCost = myIP
			}

			ElevSlice[sToi(lowestCost, ElevSlice)].Orders[buttonPress.Floor][types.ButtonMap[buttonPress.Button]] = 1

			ElevToNetCh <- ElevSlice

			time.Sleep(2 * time.Millisecond)

		case completedOrder := <-completedOrderCh: //orders executed by state machine
			fmt.Printf(" \n Floor at beginning: %v \n", ElevSlice[sToi(myIP, ElevSlice)].Floor)
			for ordNum := 0; ordNum < len(StateMachineOrderSlice); {
				if completedOrder.Floor == StateMachineOrderSlice[ordNum].Floor {
					ElevSlice[sToi(myIP, ElevSlice)].Orders[completedOrder.Floor][0] = 0
					ElevSlice[sToi(myIP, ElevSlice)].Orders[completedOrder.Floor][1] = 0
					ElevSlice[sToi(myIP, ElevSlice)].Orders[completedOrder.Floor][2] = 0
					StateMachineOrderSlice = append(StateMachineOrderSlice[:ordNum], StateMachineOrderSlice[ordNum+1:]...)
					ordNum = 0
				} else {
					ordNum++
				}
			}
			lightCh <- ElevSlice
			//updateLightsAllElevators(ElevMap, myIP, elevOnNet)
			//updateLights(ElevSlice[sToi(myIP, ElevSlice)])

			fmt.Printf("Floor: %v \n", ElevSlice[sToi(myIP, ElevSlice)].Floor)

			ElevToNetCh <- ElevSlice

		case dir := <-directionCh:
			ElevSlice[sToi(myIP, ElevSlice)].Dir = dir

		case newElevSlice := <-ElevToDistrCh:
			fmt.Printf("in new elev slice \n")

			for index := 0; index < len(newElevSlice); index++ {
				incomingID := newElevSlice[index].ID
				//HVa skjer om man får samme id to ganger?
				if incomingID == myIP {
					ElevSlice[sToi(incomingID, ElevSlice)].Orders = newElevSlice[sToi(incomingID, ElevSlice)].Orders

				} else {
					ElevSlice[sToi(incomingID, ElevSlice)] = newElevSlice[sToi(incomingID, ElevSlice)]
				}
			}

			lightCh <- ElevSlice
			//updateLightsAllElevators(ElevSlice, myIP, elevOnNet)
			//updateLights(ElevSlice[sToi(myIP, ElevSlice)])

			StateMachineOrderSlice = matrixToOrderList(ElevSlice[sToi(myIP, ElevSlice)], orderCount, StateMachineOrderSlice)
			if len(StateMachineOrderSlice) != 0 {
				orderListCh <- StateMachineOrderSlice
			}

		case elevOnNet = <-elevOnNetworkCh:
			if len(elevOnNet.Peers) == 0 {
				elevOnNet.Peers = append(elevOnNet.Peers, myIP)
			}

			if len(elevOnNet.New) > 0 {
				new := types.Elevator{}
				new.ID = elevOnNet.New
				ElevSlice = append(ElevSlice, new)
			}
			if len(ElevSlice) == 3 {
				fmt.Printf("1: %v \n", ElevSlice[0].ID)
				fmt.Printf("2: %v \n", ElevSlice[1].ID)
				fmt.Printf("3: %v \n", ElevSlice[2].ID)
			}
			//fmt.Printf("elev: %v Orders: %v \n", ElevMap[elevOnNet.Peers[ipIndex]].ID, ElevMap[elevOnNet.Peers[ipIndex]].Orders)
			if len(elevOnNet.Lost) == 1 {
				fmt.Printf("\n!!Redistributing orders!!\n")
				redistributeOrders(elevOnNet, ElevSlice, buttonEventCh, myIP)
			}
		}
	}
}

func redistributeOrders(elevOnNet peers.PeerUpdate, ElevSlice []types.Elevator,
	buttonEventCh chan<- types.ButtonEvent, ID string) {
	var ButtonPress types.ButtonEvent
	if len(elevOnNet.Lost) == 1 && elevOnNet.Peers[0] == ID {
		for floorNum := 0; floorNum < types.NFloors; floorNum++ {
			for btnNum := 0; btnNum < types.ButtonCab; btnNum++ {
				fmt.Printf("Elev lost is: %v \n", elevOnNet.Lost[0])
				fmt.Printf("Its map is: %v \n", ElevSlice[sToi(elevOnNet.Lost[0], ElevSlice)])
				fmt.Printf("ElevSlice[myID].Orders[floorNum][btnNum]: %v\n\n", ElevSlice[sToi(elevOnNet.Lost[0], ElevSlice)].Orders[floorNum][btnNum])
				if ElevSlice[sToi(elevOnNet.Lost[0], ElevSlice)].Orders[floorNum][btnNum] == 1 {
					ButtonPress.Floor = floorNum
					ButtonPress.Button = types.ButtonType(btnNum)
					ElevSlice[sToi(elevOnNet.Lost[0], ElevSlice)].Orders[floorNum][btnNum] = 0
					fmt.Printf("order sent %v \n", ButtonPress)
					buttonEventCh <- ButtonPress
				}
			}
		}
	}
}

// TimeToIdle does..
func TimeToIdle(e types.Elevator, buttonEvent types.ButtonEvent) float64 {

	duration := 0.0

	switch e.Behaviour {
	case types.Idle: //idle
		e.Dir = requests_chooseDirection(e)
		if e.Dir == types.MotorDirectionStop {
			fmt.Printf("Elev nr: %v floor %v \n", e.ID, e.Floor)
			distancePlusOne := math.Abs(float64(e.Floor-buttonEvent.Floor)) + 1
			fmt.Printf("distance plus one: %v \n", distancePlusOne)
			weight := -4 / distancePlusOne
			return weight
		}
	case types.Moving: //moving
		duration += types.TravelTime / 2
		e.Floor += int(e.Dir) //e.Dir takes values 1, 0 , -1

	case types.DoorOpen: //door open
		duration -= types.DoorOpenTime / 2
	}

	for {
		if requests_shouldStop(e) {
			e = requests_clearAtCurrentFloor(e)
			duration += types.DoorOpenTime
			e.Dir = requests_chooseDirection(e)
			if e.Dir == types.MotorDirectionStop { //MD_stop
				return duration
			}
		}
		e.Floor += int(e.Dir)
		duration += types.TravelTime
	}
}

func requests_chooseDirection(e types.Elevator) types.MotorDirection {

	var belowScore = 0.0
	var aboveScore = 0.0
	for floorNum := 0; floorNum < types.NFloors; floorNum++ {
		prio := (math.Abs(float64(e.Floor)-float64(floorNum)) + 1.0) * 1.5
		for btnNum := 0; btnNum < types.NButtons; btnNum++ {
			if floorNum < e.Floor {
				belowScore += float64(e.Orders[floorNum][btnNum]) * prio
			} else if floorNum > e.Floor {
				aboveScore += float64(e.Orders[floorNum][btnNum]) * prio
			}
		}
	}
	if belowScore == 0 && aboveScore == 0 {
		return types.MotorDirectionStop
	} else if belowScore < aboveScore {
		return types.MotorDirectionUp
	} else {
		return types.MotorDirectionDown
	}
}

func requests_shouldStop(e types.Elevator) bool {

	switch e.Dir {
	case types.MotorDirectionDown:
		return (e.Orders[e.Floor][2] == 1 || e.Orders[e.Floor][1] == 1 || e.Floor == 0)
	case types.MotorDirectionUp:
		return (e.Orders[e.Floor][2] == 1 || e.Orders[e.Floor][0] == 1 || e.Floor == 3)
	case types.MotorDirectionStop:
		return true
	default:
		return false
	}
}

func requests_clearAtCurrentFloor(e_old types.Elevator) types.Elevator {

	e := e_old
	var btn types.ButtonType

	for ; btn < types.NButtons; btn = btn + 1 {
		if e.Orders[e.Floor][btn] == 1 {
			e.Orders[e.Floor][btn] = 0
		}
	}
	return e
}

func isDuplicate(b types.ButtonEvent, ElevSlice []types.Elevator, p peers.PeerUpdate, myIP string) bool {

	btnInt := types.ButtonMap[b.Button]

	if btnInt == 2 {
		return (ElevSlice[sToi(myIP, ElevSlice)].Orders[b.Floor][btnInt] == 1)

	} else {
		for elevIndex := 0; elevIndex < len(ElevSlice); elevIndex++ {
			if ElevSlice[sToi(myIP, ElevSlice)].Orders[b.Floor][btnInt] == 1 {
				return true
			}
		}
	}
	return false
}

func matrixToOrderList(e types.Elevator, orderCount int, list []types.SingleOrder) []types.SingleOrder {

	for floorNum := 0; floorNum < types.NFloors; floorNum++ {
		for btnNum := 0; btnNum < types.NButtons; btnNum++ {
			if e.Orders[floorNum][btnNum] == 1 {
				tempOrder := types.SingleOrder{" ", types.ButtonType(btnNum), floorNum, false}

				if !isInOrderList(tempOrder, list) {
					list = append(list, tempOrder)
					orderCount++
				}
			}
		}
	}
	return list
}

func isInOrderList(s types.SingleOrder, list []types.SingleOrder) bool {

	for index := 0; index < len(list); index++ {
		if (s.Floor == list[index].Floor) && (s.Button == list[index].Button) {
			return true
		}
	}
	return false
}

func updateMatrixAndAppend(eOld types.Elevator, eNew types.Elevator, OrderList []types.SingleOrder, orderCount int) ([]types.SingleOrder, types.Elevator) {

	for floorNum := 0; floorNum < types.NFloors; floorNum++ {
		for btnNum := 0; btnNum < types.NButtons; btnNum++ {
			if (eNew.Orders[floorNum][types.ButtonType(btnNum)] == 1) && (eOld.Orders[floorNum][types.ButtonType(btnNum)] != 1) {
				OrderList = append(OrderList, types.SingleOrder{ID: strconv.Itoa(orderCount), Button: types.ButtonType(btnNum), Floor: floorNum, MsgConf: false})
				orderCount++
			}
			eNew.Orders[floorNum][types.ButtonType(btnNum)] = (eOld.Orders[floorNum][types.ButtonType(btnNum)] | eNew.Orders[floorNum][types.ButtonType(btnNum)])
		}
	}
	return OrderList, eNew

}

//ChanUpdateLight Inspired by Per Kjelsvik
func ChanUpdateLight(lightCh chan []types.Elevator, ID string) {

	for {
		elevs := <-lightCh
		sToi(ID, elevs)

		for floorNum := 0; floorNum < types.NFloors; floorNum++ {
			driver.SetButtonLamp(types.ButtonCab, floorNum, elevs[sToi(ID, elevs)].Orders[floorNum][types.ButtonCab] == 1)
		}

		exists := false
		for floor := 0; floor < types.NFloors; floor++ {
			for btn := 0; btn < types.NButtons-1; btn++ {
				exists = false
				for elevator := 0; elevator < len(elevs); elevator++ {
					if elevs[elevator].Orders[floor][btn] == 1 {
						exists = true
					}
				}
				driver.SetButtonLamp(types.ButtonType(btn), floor, exists)
			}
		}
	}
}

func updateLights(e types.Elevator) {
	for floorNum := 0; floorNum < types.NFloors; floorNum++ {
		for btnNum := 0; btnNum < types.NButtons; btnNum++ {
			driver.SetButtonLamp(types.ButtonType(btnNum), floorNum, e.Orders[floorNum][btnNum] == 1)
		}
	}
}

// ElevatorMotorError turns of the network when a motor error is detected. It only returns true when a motor error is present and the button pressed is a Cab orders
func ElevatorMotorError(buttonPress types.ButtonEvent, motorErrorCh chan bool, turnOfNetworkCh chan bool) bool {
	select {
	case MotorError = <-motorErrorCh:
		turnOfNetworkCh <- true
		return (buttonPress.Button != types.ButtonCab) && MotorError
	default:
		return (buttonPress.Button != types.ButtonCab) && MotorError
	}
}

func sToi(ElevatorID string, elevSlice []types.Elevator) int {

	for index := 0; index < len(elevSlice); index++ {
		if ElevatorID == elevSlice[index].ID {
			return index
		}
	}
	return 0
}
