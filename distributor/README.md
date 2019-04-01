# Module distributor

- repsonsible for the reading and distribution of new orders \
&nbsp;&nbsp;&nbsp;&nbsp;- the redistribution is achieved thorugh calculating time to ide for all elevators \
&nbsp;&nbsp;&nbsp;&nbsp;- the elevator eith the lowest time to idle is assigned the order
- updates the lights
- communicates with all other modules through channles <br/>
- saves the state and orders of all other elevators in a slice of type types.Elevator <br/>
- handles any fault cases that might occur <br/> <br/> <br/>


*type Elevator struct { \
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  Floor     int\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	Dir       MotorDirection\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	Orders    [NFloors][NButtons]int\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	Behaviour ElevatorBehaviour\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;	ID        string \
}* <br/> <br/>



Floors | Hall Up   | Hall Down  |     Cab    |
----------- | ---------- | ---------- | ----------
Floor 3     | - |- | -
Floor 2     | - |- | -
Floor 1     | - | - | -
Floor 0     | - | - | -
