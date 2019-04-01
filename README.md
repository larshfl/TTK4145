# Gruppe 1 - Sanntidsprogrammering TTK4145 - Vår 2019

## Qucik start 
### Run elevator on the hardware used at "sanntidslabben"
```ruby 
$ ./startWithSafety.sh <port number>
```

## Intro


## System descripton? Should we have the flow diagram here?


## Code inspiration
We have used several modules written by @klasbo.

The driver module is mainly based on https://github.com/TTK4145/driver-go/tree/master/elevio<br />
Changes: structs and consts have been renamed and moved to types


The network module uses conn, bcast, localip and peers. These are almost an exact copy of https://github.com/TTK4145/Network-go/tree/master/network<br />
Changes:

The statemachine is inspired by Anders Petersens (@klasbo) lecture https://www.youtube.com/watch?v=K6YoNYNC7o4&t=1646s where he used a statemachine written in go as a visual aid. The elevator system design differs, but the structure has similarities.

The distributor module uses TimeToIdle as well as all functions called inside this function. These functions are written by (@klasbo) and taken from https://github.com/TTK4145/Project-resources. The functions have been translated into Go syntax and some have been slightly modified.


## Legg også inn om vi har lånt kode fra nettet - Kan få problem med plagiat om vi har "lånt" kode uten å gi credit 

## Skriv readme til kvar av modulane

## Legg in bilde av flyskjema for tolatsystemet og for nettwork



## Abbreviations
___
Ch   &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;   Channel <br />
Msg  &nbsp;&nbsp;&nbsp;&nbsp;   Message <br />
Num  &nbsp;&nbsp;&nbsp;&nbsp;   Number <br />
___




## Testing
The development of this elevator system required lots of immaculate testing. We used the simulator written by @Klasbon as the main test environment. The simulator is available in this repository [SimElevator](https://github.com/TTK4145-students-2019/project-group-1/blob/master/SimElevatorServer). 

In order to make testing with multiples computers as smooth as possible we're using SSH to get access multiple computers. After a connection is established we run a shell script. That checks if a folder Gruppe1 exists on the Desktop and if the repo exisits in that folder the repo will be deleted and the latets version will be cloned


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

