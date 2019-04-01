# Module network

The network module handles peers and message communication between the elevators. The module is based around the Messagestruct from types, the workflow can be seen illustrated in the flowchart below.

![alt text](https://i.imgur.com/LHpxfGH.png)

In the implemented code the tx-queue is implemented as a buffered channel from distributor to network, since channels act as FIFO-queues, this saves us quite a bit of code. The left decision loop of the flowchart is implemented through the function "txMsgHandler" where a for select used is used to check the channel from distributor for incoming messages that should be sent out and confirmed by all elevators available in peers. To control peer and message confirmation, two global slices are used with accompanying mutexes. It is possible to implement a message based rather than a lock based solution, but a lock based solution provided readability for the code. The confirmation check is implemented through a simple subsetfunction which checks if current peer IDs are a subset of current confirmation IDs.

##Message struct

'''go
type Message struct {
	ElevID            string
	Content           []Elevator
	ConfirmedMsgOwner string
	MsgID             int
}
'''