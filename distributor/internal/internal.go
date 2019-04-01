package internal

import (
	"../../types"
	"strconv"
)

// RedistributeOrders distributes order of the disconnected elevators
func RedistributeOrders(ElevSlice []types.Elevator,	buttonEventCh chan<- types.ButtonEvent, intID int) {

	var ButtonPress types.ButtonEvent
	for floorNum := 0; floorNum < types.NumFloors; floorNum++ {
		for btnNum := 0; btnNum < types.ButtonCab; btnNum++ {
			if ElevSlice[intID].Orders[floorNum][btnNum] == 1 {
				ButtonPress.Floor = floorNum
				ButtonPress.Button = types.ButtonType(btnNum)
				ElevSlice[intID].Orders[floorNum][btnNum] = 0
				buttonEventCh <-ButtonPress
			}
		}
	}
}

// IsDuplicate checks for duplicate orders
func IsDuplicate(b types.ButtonEvent, ElevSlice []types.Elevator, ID int) bool {

	btnInt := types.ButtonMap[b.Button]

	if btnInt == 2 {
		return (ElevSlice[ID].Orders[b.Floor][btnInt] == 1)
	}
	for elevIndex := 0; elevIndex < types.NumElevators; elevIndex++ {
		if ElevSlice[ID].Orders[b.Floor][btnInt] == 1 {
			return true
		}
	}
	return false
}

// MatrixToOrderList convert the order matrix to the orderlist for the statemachine
func MatrixToOrderList(e types.Elevator, list []types.SingleOrder) []types.SingleOrder {

	for floorNum := 0; floorNum < types.NumFloors; floorNum++ {
		for btnNum := 0; btnNum < types.NumButtons; btnNum++ {
			if e.Orders[floorNum][btnNum] == 1 {
				tempOrder := types.SingleOrder{Button: types.ButtonType(btnNum), Floor: floorNum}
				if !IsInOrderList(tempOrder, list) {
					list = append(list, tempOrder)
				}
			}
		}
	}
	return list
}

// IsInOrderList checks if the order is in the oderlist
func IsInOrderList(s types.SingleOrder, list []types.SingleOrder) bool {

	for index := 0; index < len(list); index++ {
		if (s.Floor == list[index].Floor) && (s.Button == list[index].Button) {
			return true
		}
	}
	return false
}

// ElevSliceInit initializes state and beahavior of the Elevslice
func ElevSliceInit(elevSlice []types.Elevator, ID int) []types.Elevator {

	for elev := 0; elev < types.NumElevators; elev++ {
		if elev == ID {
			elevSlice[elev].Behaviour = types.Idle
			elevSlice[elev].ID = strconv.Itoa(elev)
		} else {
			elevSlice[elev].Behaviour = types.Undefined
			elevSlice[elev].ID = strconv.Itoa(elev)
		}
	}
	return elevSlice
}