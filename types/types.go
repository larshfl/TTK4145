package types

// const (
// 	Idle elevatorBehavior = iota
// 	Moving
// 	DoorOpen
// )

// Elevator states
const (
	Idle     = 0
	Moving   = 1
	DoorOpen = 2
)

// MotorDirection is the type for the three motor directions
type MotorDirection int

// Motor Directions
const (
	MotorDirectionUp   MotorDirection = 1
	MotorDirectionDown                = -1
	MotorDirectionStop                = 0
)

// ButtonType is the type
type ButtonType int

//Buttons
const (
	ButtonHallUp   ButtonType = 0
	ButtonHallDown            = 1
	ButtonCab                 = 2
)

// ElevatorBehaviour is the type for elevator behaviour
type ElevatorBehaviour int

// Elevator struct
type Elevator struct {
	Floor     int
	Dir       MotorDirection
	Orders    [NFloors][NButtons]int // |Up	|Down	|Cab	|
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

// TravelTime for the elevator between two floors
const TravelTime = 3.0 // reisetid mellom etasjer
// DoorOpenTime for the elevator
const DoorOpenTime = 3.0

// NButtons is the number of different buttontypes
const NButtons = 3

// NFloors is the number of floors
const NFloors = 4

// SingleOrder is the ID, buttonType and Floor
type SingleOrder struct {
	ID      string
	Button  ButtonType
	Floor   int
	MsgConf bool
}

// ButtonEvent is the Floor and buttontype
type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// ButtonMap maps the ButtonType to an int
var ButtonMap = map[ButtonType]int{
	ButtonHallUp: 0,
	ButtonHallDown: 1,
	ButtonCab: 2,
}
