package cost

import (
	"math"
	
	"../../types"
)


// TimeToIdle calculates cost of orders
func TimeToIdle(e types.Elevator, buttonEvent types.ButtonEvent) float64 {

	duration := 0.0

	switch e.Behaviour {
	case types.Idle:
		e.Dir = chooseDirection(e)
		if e.Dir == types.MotorDirectionStop {
			distancePlusOne := math.Abs(float64(e.Floor-buttonEvent.Floor)) + 1
			weight := -4 / distancePlusOne
			return weight
		}
	case types.Moving:
		duration += types.TravelTime / 2
		e.Floor += int(e.Dir)

	case types.DoorOpen:
		duration -= types.DoorOpenTime / 2
	}

	for {
		if shouldStop(e) {
			e = clearAtCurrentFloor(e)
			duration += types.DoorOpenTime
			e.Dir = chooseDirection(e)
			if e.Dir == types.MotorDirectionStop {
				return duration
			}
		}
		e.Floor += int(e.Dir)
		duration += types.TravelTime
	}
}

func chooseDirection(e types.Elevator) types.MotorDirection {

	var belowScore = 0.0
	var aboveScore = 0.0
	for floorNum := 0; floorNum < types.NumFloors; floorNum++ {
		prio := (math.Abs(float64(e.Floor)-float64(floorNum)) + 1.0) * 1.5
		for btnNum := 0; btnNum < types.NumButtons; btnNum++ {
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

func shouldStop(e types.Elevator) bool {

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

func clearAtCurrentFloor(eOld types.Elevator) types.Elevator {

	e := eOld

	for btn := types.ButtonType(0); btn < types.NumButtons; btn++ {
		if e.Orders[e.Floor][btn] == 1 {
			e.Orders[e.Floor][btn] = 0
		}
	}
	return e
}
