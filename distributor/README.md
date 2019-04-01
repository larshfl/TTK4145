# Module distributor

- repsonsible for the reading and distribution of new orders
- updates the lights
- communicates with all other modules through channles
- saves the state and orders of all other elevators in a slice of type types.Elevator


*type Elevator struct { \
	Floor     int\
	Dir       MotorDirection\
	Orders    [NFloors][NButtons]int\
	Behaviour ElevatorBehaviour\
	ID        string //bruke IP som ID?\
}*\



Floors | Hall Up   | Hall Down  |     Cab    |
----------- | ---------- | ---------- | ----------
Floor 3     | - |- | -
Floor 2     | - |- | -
Floor 1     | - | - | -
Floor 0     | - | - | -
