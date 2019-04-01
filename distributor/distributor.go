package distributor

import (
	"math"
	"strconv"
	"time"

	"../network/peers"
	"../types"
	
	"./cost"
	"./internal"
)
// Distributor calculates cost and distributes orders
func Distributor(
	currentFloorCh 		chan int,
	buttonEventCh 		chan types.ButtonEvent,
	elevOnNetworkCh 	chan peers.PeerUpdate,
	completedOrderCh 	chan types.SingleOrder,
	directionCh 		chan types.MotorDirection,
	motorErrorCh 		chan bool,
	ElevToNetCh 		chan []types.Elevator,
	networkEnableCh 	chan bool,
	orderListCh 		chan []types.SingleOrder,
	singleOrderCh 		chan types.SingleOrder,
	ElevToDistrCh		chan []types.Elevator,
	lightCh 			chan []types.Elevator,
	ID 					int) {

	ElevSlice := make([]types.Elevator, types.NumElevators)
	ElevSlice = internal.ElevSliceInit(ElevSlice, ID)
	StateMachineOrderSlice := make([]types.SingleOrder, 0)
	var elevOnNet peers.PeerUpdate
	motorError := false

	for {
		select {
		case floor := <-currentFloorCh:
			motorError = false
			networkEnableCh <- true
			ElevSlice[ID].Floor = floor

		case motorError := <-motorErrorCh:
			networkEnableCh <-!motorError

		case buttonPress := <-buttonEventCh:
			if motorError && (buttonPress.Button != types.ButtonCab) {
				break
			}
			if internal.IsDuplicate(buttonPress, ElevSlice, ID) {
				break
			}

			lowestCost := ID
			min := math.Inf(1)

			for ipIndex := 0; ipIndex < types.NumElevators; ipIndex++ {
				if ElevSlice[ipIndex].Behaviour != types.Undefined {
					cost := cost.TimeToIdle(ElevSlice[ipIndex], buttonPress)
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
			ElevToNetCh <-ElevSlice

			// To make sure two buttonpresses aren't evaulated before the first change has taken effect 
			time.Sleep(2 * time.Millisecond)

		case completedOrder := <-completedOrderCh:
			for ordNum := 0; ordNum < len(StateMachineOrderSlice); {
				if completedOrder.Floor == StateMachineOrderSlice[ordNum].Floor {
					for btn := 0; btn < types.NumButtons; btn++ {
						ElevSlice[ID].Orders[completedOrder.Floor][btn] = 0
					}
					StateMachineOrderSlice = append(StateMachineOrderSlice[:ordNum], StateMachineOrderSlice[ordNum+1:]...)
					ordNum = 0
				} else {
					ordNum++
				}
			}

			ElevToNetCh <-ElevSlice

		case dir := <-directionCh:
			ElevSlice[ID].Dir = dir

		case newElevSlice := <-ElevToDistrCh:
			for index := 0; index < types.NumElevators; index++ {
				incomingID, _ := strconv.Atoi(newElevSlice[index].ID)
				if incomingID == ID {
					ElevSlice[ID].Orders = newElevSlice[ID].Orders
				} else {
					ElevSlice[incomingID] = newElevSlice[incomingID]
				}
			}
			lightCh <-ElevSlice
			StateMachineOrderSlice = internal.MatrixToOrderList(ElevSlice[ID], StateMachineOrderSlice)
			if len(StateMachineOrderSlice) != 0 {
				orderListCh <- StateMachineOrderSlice
			}

		case elevOnNet = <-elevOnNetworkCh:
			for IDindex := 0; IDindex < types.NumElevators; IDindex++ {
				if ElevSlice[IDindex].ID == elevOnNet.New {
					ElevSlice[IDindex].Behaviour = types.Idle
				}
			}

			for index := 0; index < len(elevOnNet.Lost); index++ {
				intID, _ := strconv.Atoi(elevOnNet.Lost[index])
				ElevSlice[intID].Behaviour = types.Undefined
				internal.RedistributeOrders(ElevSlice, buttonEventCh, intID)
			}
		}
	}
}
