# Sanntidsprogrammering TTK4145 - Vår 2019

## Quick start 

### Clone the repo
```ruby 
$ ./clone.sh [path: $HOME/Desktop]
```
This script check if the folder Gruppe1 is a directory in $HOME/Desktop. If it don't exists it will make the directory $HOME/Desktop/Gruppe1 and clone the repositort to this directory. If $HOME/Desktop/Gruppe1 exists it will check for a repository and delete it befor cloning the latest version available on GitHub

### Run elevator on the hardware used at "sanntidslabben"

```ruby 
$ ./gruppe1.sh <port number>
```
This starts the Elevator server used to communicate with the elevators at "sanntidslabben". After it starts the Elevator server it starts the main program with the ID provided by the user. In this case 0, 1 or 2. 

### Run elevators on the simulator
```ruby 
$ ./elevatorShell.sh
```
This starts three simulators on the default ports 10001, 10002 and 10003. And starts three instanses of the main 
program with the same ports and default ID's 0, 1 and 2.


## Project description
TTK4145 Real-Time Programming has a semester project that in short are; 
Create software for controlling `n` elevators working in parallel across `m` floors.

The project have several requirements. The most important are.


## Elevator project
### Description
In this project, we had to create software for controlling `n` elevators working in parallel across `m` floors. There were some main requirments, here summarised in bullet points: 

  - **No orders are lost** 
  - **Multiple elevators should be more efficient than one** 
  - **An Inidividual elevator should behave sensibly and efficiently**
  - **The lights should function as expected**
  
In the project, we start with `1 <= n <= 3` elevators, and `m == 4` floors. However, we should avoid hard-coding these values, and we aimed to write the project where adding a floor or an elevator required minimal work. The system will however _not_ be tested for `n > 3` or `m != 4`. There are also some unspecified behaviours we had to decide for ourselves: 

  - **Which orders are cleared when stopping at a floor**
  - **How the elevator behaves when it cannot connect to the network during initialization**
  - **How the hall (call up, call down) buttons work when the elevator is disconnected from the network**
 
Lastly, there were some permitted assumptions: 

  - **At least one elevator is always working normally**
  - **No multiple simultaneous errors: Only one error happens at a time, but the system must still return to a fail-safe state after this error**
  - **No network partitioning: Situations where there are multiple sets of two or more elevators with no connection between them can be ignored**
  - **Stop button and obstruction switch are disabled**

For full details on each point, the driver files, or the full specs of the project: head over to [`TTK4145`](https://github.com/TTK4145/Project#elevator-project) (the description/project might have changed as the course is held every spring)


### Our solution 
We wrote the elevator system in Golang and devided the system into four main modules; statemachine, distributor, driver and nettwork. The four modules and its interface of is shown in the flow diagram below. We're mainly using channels as a way of communcation between the modules with the exception of two functioncalls to the driver module in order to toggle the lights. The picture below also shows the elevator struct that is passed on the network and from network to distributor 
![Flow diagram](https://i.imgur.com/fSjjoZ9.png)



### State Machine
The State Machine controls the elevator and executes the orders sent to it by the distributor module. The states are only evaluated and changed upon reception of data over the channels. Most of the interaction with the elevator hardware is done thorugh this module. 

### Distributor
Most of the program is controled through the distributor. This module receives and assigns new orders to one of the eleavtors present on the network and sends these orders either to the local state machine or to the network module which forwards it to  one of the other elevators. The distributor module is furthermore tasked with handeling the fault cases that migh occur. 

### Driver
This is the hardware abstracton module which communicates with the hardware over TCP.

### Network
The Network module implements an acknowledge LAN communication service based upon the UDP protocol.
![Network Handshake](https://i.imgur.com/ubruIMN.png)







## Code inspiration
We have used several modules written by @klasbo.

The driver module is mainly based on https://github.com/TTK4145/driver-go/tree/master/elevio<br />
Changes: structs and consts have been renamed and moved to types

The network module uses conn, bcast and peers. These are almost an exact copy of https://github.com/TTK4145/Network-go/tree/master/network<br />
Changes:

The statemachine is inspired by Anders Petersens (@klasbo) lecture https://www.youtube.com/watch?v=K6YoNYNC7o4&t=1646s where he used a statemachine written in go as a visual aid. The elevator system design differs, but the structure have similarities.

The distributor module uses TimeToIdle as well as all functions called inside this function. These functions are written by (@klasbo) and taken from https://github.com/TTK4145/Project-resources. The functions have been translated into Go syntax and some have been slightly modified.


## Abbreviations
___
Ch   &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;   Channel <br />
Msg  &nbsp;&nbsp;&nbsp;&nbsp;   Message <br />
Num  &nbsp;&nbsp;&nbsp;&nbsp;   Number <br />
Elev &nbsp;&nbsp;&nbsp;&nbsp;   Elevator <br />
___




## Testing
The development of this elevator system required lots of immaculate testing. We used the simulator written by @Klasbo as the main test environment. The simulator is available in this repository [SimElevator](https://github.com/TTK4145-students-2019/project-group-1/blob/master/SimElevatorServer). 

In order to make testing with multiples computers as smooth as possible we're using SSH to get access multiple computers. After a connection is established we run a shell script. That checks if a folder Gruppe1 exists on the Desktop and if the repo exisits in that folder the repo will be deleted and the latets version will be cloned.

For testing with packet loss did we make a shell script that sets the packet loss to 20 %. 
### Run packet loss script
```ruby 
$ packet_loss.sh
```
### To turn off packet loss
```ruby 
$ sudo iptables -L
```




## Connect with SSH

### Preparations on linux based server
```ruby 
$ sudo apt update
$ sudo apt install openssh-server
```
**Check if SSH service is running**
```ruby 
$ sudo systemctl status ssh
```

**Start or stop SSH service**
```ruby 
$ sudo systemctl start ssh
$ sudo systemctl stop ssh
```

**To obtain the servers ip adress**
```ruby 
$ ip a
```
Use the inet ip adress. See the picture for assitance 
[ip](https://i.imgur.com/McevWcV.png)

### Connect to server through SSH
```ruby 
$ server_username@ip -Y 
```
The -Y flag enables display on the server

### Troubleshooting 
Make sure that server has open policy. Check if policys are listed as ACCPETED by running
```ruby 
$ sudo iptables -L
```
Change the status with
```ruby 
$ iptables --policy INPUT ACCEPT
```
To change it back use
```ruby 
$ iptables --policy INPUT DROP
```
