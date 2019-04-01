package types

// Elevator states
const (
	Undefined 			= iota - 1
	Idle      
	Moving    
	DoorOpen  
)

// MotorDirection is the type for the three motor directions
type MotorDirection int

// Motor direction
const (
	MotorDirectionDown  = iota - 1
	MotorDirectionStop                
	MotorDirectionUp
)

// ButtonType is the type
type ButtonType int

//Buttons
const (
	ButtonHallUp    	= iota
	ButtonHallDown            
	ButtonCab                 
)

// ElevatorBehaviour is the type for elevator behaviour
type ElevatorBehaviour int

// Elevator struct
type Elevator struct {
	Floor     int
	Dir       MotorDirection
	Orders    [NumFloors][NumButtons]int // |Up	|Down	|Cab	|
	Behaviour ElevatorBehaviour
	ID        string //bruke IP som ID?
}

//Message struct contains all info needed for communication.
type Message struct {
	ElevID            string
	Content           []Elevator
	ConfirmedMsgOwner string
	MsgID             int
}

// Constants
const (
	TravelTime					= 3.0
	DoorOpenTime 				= 3.0
	NumElevators				= 3
	NumButtons 					= 3
	NumFloors 					= 4
	MaxTimeBeforeMotorError 	= 5
)

// SingleOrder is the ID, buttonType and Floor
type SingleOrder struct {
	Button  ButtonType
	Floor   int
}

// ButtonEvent is the Floor and buttontype
type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// ButtonMap maps the ButtonType to an int
var ButtonMap = map[ButtonType]int{
	ButtonHallUp:   0,
	ButtonHallDown: 1,
	ButtonCab:      2,
}
