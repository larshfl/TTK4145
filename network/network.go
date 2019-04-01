package network

import (
	"../types"
	"./bcast"
	"./peers"
	"./communication"
)

// Network module
func Network(
	fromDistributorCh 	chan []types.Elevator,
	toDistributorCh 	chan []types.Elevator,
	peerTxEnableCh 		chan bool,
	peerDistributorCh 	chan peers.PeerUpdate,
	portNum 			string,
	ID 					string) {
		

	peerElevPort 		:= 22032
	elevMsgPort 		:= 19032

	peerCh 				:= make(chan peers.PeerUpdate)
	txCh 				:= make(chan types.Message)
	rxCh 				:= make(chan types.Message)
	msgCh 				:= make(chan types.Message, 10)
	receivedUpdateCh 	:= make(chan types.Message)
	msgID 				:= 1

	go peers.Transmitter(peerElevPort, ID, peerTxEnableCh)
	go peers.Receiver(peerElevPort, peerCh)
	go bcast.Transmitter(elevMsgPort, txCh)
	go bcast.Receiver(elevMsgPort, rxCh)
	go communication.TxMsgHandler(txCh, msgCh, ID, receivedUpdateCh)
	go communication.RxMsgHandler(rxCh, peerCh, ID, peerDistributorCh, toDistributorCh, txCh)
	for {
		select {
		case msg := <-fromDistributorCh:
			newMsg := communication.ComposeMsg(ID, msg, msgID)
			msgCh <- newMsg
			msgID++
		case newUpdate := <-receivedUpdateCh:
			toDistributorCh <- newUpdate.Content
		}
	}
}