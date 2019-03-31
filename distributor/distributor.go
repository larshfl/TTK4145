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

	ElevSlice := make([]types.Elevator, types.NElevators)
	myID := <-myIDCh
	ID, _ := strconv.Atoi(myID)
	ElevSlice = elevSliceInit(ElevSlice, ID)
	StateMachineOrderSlice := make([]types.SingleOrder, 0)
	orderCount := 0
	var elevOnNet peers.PeerUpdate

	for {
		select {
		case floor := <-currentFloorCh:
			MotorError = false
			ElevSlice[ID].Floor = floor

		case buttonPress := <-buttonEventCh:
			if ElevatorMotorError(buttonPress, motorErrorCh, turnOfNetworkCh) {
				break
			}
			if isDuplicate(buttonPress, ElevSlice, ID) {
				break
			}

			lowestCost := ID
			min := math.Inf(1)

			for ipIndex := 0; ipIndex < types.NElevators; ipIndex++ {
				if ElevSlice[ipIndex].Behaviour != types.Undefined {
					cost := TimeToIdle(ElevSlice[ipIndex], buttonPress)
					if min > cost {
						min = cost
						lowestCost, _ = strconv.Atoi(ElevSlice[ipIndex].ID)
					}
				}
			}

			if types.ButtonCab == buttonPress.Button {
				lowestCost = ID
			}
			ElevSlice[lowestCost].Orders[buttonPress.Floor][types.ButtonMap[buttonPress.Button]] = 1
			ElevToNetCh <- ElevSlice
			time.Sleep(2 * time.Millisecond)

		case completedOrder := <-completedOrderCh: //orders executed by state machine
			fmt.Printf("Completed order received \n")
			for ordNum := 0; ordNum < len(StateMachineOrderSlice); {
				if completedOrder.Floor == StateMachineOrderSlice[ordNum].Floor {
					ElevSlice[ID].Orders[completedOrder.Floor][0] = 0
					ElevSlice[ID].Orders[completedOrder.Floor][1] = 0
					ElevSlice[ID].Orders[completedOrder.Floor][2] = 0
					StateMachineOrderSlice = append(StateMachineOrderSlice[:ordNum], StateMachineOrderSlice[ordNum+1:]...)
					ordNum = 0
				} else {
					ordNum++
				}
			}

			fmt.Printf("Elev Slice = %v \n", ElevSlice[ID].Orders)
			lightCh <- ElevSlice
			ElevToNetCh <- ElevSlice

		case dir := <-directionCh:
			ElevSlice[ID].Dir = dir

		case newElevSlice := <-ElevToDistrCh:
			for index := 0; index < types.NElevators; index++ {
				incomingID, _ := strconv.Atoi(newElevSlice[index].ID)
				if incomingID == ID {
					ElevSlice[ID].Orders = newElevSlice[ID].Orders

				} else {
					ElevSlice[incomingID] = newElevSlice[incomingID]
				}
			}

			lightCh <- ElevSlice
			StateMachineOrderSlice = matrixToOrderList(ElevSlice[ID], orderCount, StateMachineOrderSlice)
			fmt.Printf("StateMachine order slice %v \n", StateMachineOrderSlice)
			if len(StateMachineOrderSlice) != 0 {
				orderListCh <- StateMachineOrderSlice
			}

		case elevOnNet = <-elevOnNetworkCh:
			// if len(elevOnNet.Peers) == 0 {
			// 	elevOnNet.Peers = append(elevOnNet.Peers, ID)
			// }

			for IDindex := 0; IDindex < types.NElevators; IDindex++ {
				if ElevSlice[IDindex].ID == elevOnNet.New {
					ElevSlice[IDindex].Behaviour = types.Idle
				}
			}

			for index := 0; index < len(elevOnNet.Lost); index++ {
				intID, _ := strconv.Atoi(elevOnNet.Lost[index])
				ElevSlice[intID].Behaviour = types.Undefined
				fmt.Printf("\n!!Redistributing orders!!\n")
				redistributeOrders(elevOnNet, ElevSlice, buttonEventCh, intID)
			}
		}
	}
}

func redistributeOrders(elevOnNet peers.PeerUpdate, ElevSlice []types.Elevator,
	buttonEventCh chan<- types.ButtonEvent, intID int) {
	var ButtonPress types.ButtonEvent
	for floorNum := 0; floorNum < types.NFloors; floorNum++ {
		for btnNum := 0; btnNum < types.ButtonCab; btnNum++ {

			if ElevSlice[intID].Orders[floorNum][btnNum] == 1 {
				ButtonPress.Floor = floorNum
				ButtonPress.Button = types.ButtonType(btnNum)
				ElevSlice[intID].Orders[floorNum][btnNum] = 0
				buttonEventCh <- ButtonPress
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

			distancePlusOne := math.Abs(float64(e.Floor-buttonEvent.Floor)) + 1

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

func isDuplicate(b types.ButtonEvent, ElevSlice []types.Elevator, ID int) bool {

	btnInt := types.ButtonMap[b.Button]

	if btnInt == 2 {
		return (ElevSlice[ID].Orders[b.Floor][btnInt] == 1)

	} else {
		for elevIndex := 0; elevIndex < types.NElevators; elevIndex++ {
			if ElevSlice[ID].Orders[b.Floor][btnInt] == 1 {
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

func elevSliceInit(elevSlice []types.Elevator, ID int) []types.Elevator {
	for elev := 0; elev < types.NElevators; elev++ {
		intID, _ := strconv.Atoi(elevSlice[elev].ID)
		if intID == ID {
			elevSlice[elev].Behaviour = types.Idle
		} else {
			elevSlice[elev].Behaviour = types.Undefined
		}
	}
	return elevSlice
}
