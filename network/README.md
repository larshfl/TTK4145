# Module network

The network module handles peers and message communication between the elevators. The module is based around the Messagestruct from types, the workflow can be seen illustrated in the flowchart below.

![alt text](https://i.imgur.com/LHpxfGH.png)

In the implemented code the tx-queue is implemented as a buffered channel from distributor to network, since channels act as FIFO-queues, this saves us quite a bit of code. The left decision loop of the flowchart is implemented through the function "txMsgHandler" where a for select used is used to check the channel from distributor for incoming messages that should be sent out and confirmed by all elevators available in peers. To control peer and message confirmation, two global slices are used with accompanying mutexes. It is possible to implement a message based rather than a lock based solution, but a lock based solution provided readability for the code. The confirmation check is implemented through a simple subsetfunction which checks if current peer IDs are a subset of current confirmation IDs.

# Message struct #

The Message struct contains everything required for communication described by the above flowchart. It contains an ElevID which is the unique ID of the elevator, content, which contains message that should be sent to distributor module, a confirmedMsgOwner string which is the unique ID of an elevator that should receive a confirmation, and also a MsgID which is incremented, so that older messages are not sent to the distributor.

```go
type Message struct {
	ElevID            string
	Content           []Elevator
	ConfirmedMsgOwner string
	MsgID             int
}
```

If an elevator receives a message, it is handled in the rxMsgHandler, which handles peer updates from the peer module and incoming message from the UDPreceiver. If a peer update is received it updates the distributor and updates the shared list of peers. If a message is received from the UDPreceiver it checks how it should handle the message by checking confirmedMsgOwner, MsgID and sender-ID. If it is a "content-message" it sends a confirmation message and updates distributor if the message is newer thatn the previously received from the senders ID.

The transmitting, reception, and formatting of messages is done in the bcast module, which is written by @klasbo and @kjetilkjeka. They have also written the peers-module used heavily in network. 
