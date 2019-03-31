package setup

import (
	"../driver"
	"../types"
	"fmt"
)

// Init - initialize elevator
func Init(portNum string) {
	driver.Init(fmt.Sprintf("localhost:%v",portNum), types.NFloors) 

	// Turn of all Button Lamps
	for floor := 0; floor < types.NFloors; floor++ {
		for button := types.ButtonType(0); button < 3; button++ {
			driver.SetButtonLamp(button, floor, false)
		}
	}

	// Turn of stop lamp
	driver.SetStopLamp(false)
}
