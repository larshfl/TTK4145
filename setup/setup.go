package setup

import (
	"../driver"
	"../types"
	"fmt"
)

// Init - initialize elevator
func Init(portNum string) {
	driver.Init(fmt.Sprintf("localhost:%v",portNum), types.NumFloors) 

	for floor := 0; floor < types.NumFloors; floor++ {
		for button := types.ButtonType(0); button < 3; button++ {
			driver.SetButtonLamp(button, floor, false)
		}
	}

	driver.SetStopLamp(false)
}
