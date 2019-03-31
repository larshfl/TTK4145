package network

import (
	//"strconv"
	//"fmt"
	"sync"
	"time"

	"../types"
	"./bcast"
	"./peers"
	//"./localip"
	//"fmt"
	//"strconv"
)

var currentConfirmations []string
var currentPeers []string
var confirmationMutex = &sync.Mutex{}
var peerMutex = &sync.Mutex{}

func composeMsg(localip string, content []types.Elevator, msgID int) types.Message {
	var newMsg types.Message
	newMsg.ElevID = localip
	newMsg.MsgID = msgID
	newMsg.Content = content
	return newMsg
}

func sendConfirmation(rxElevID string, thisElevID string, msgTxCh chan types.Message) {
	var confirmationMsg types.Message
	confirmationMsg.Content = nil
	confirmationMsg.ConfirmedMsgOwner = rxElevID
	confirmationMsg.ElevID = thisElevID
	msgTxCh <- confirmationMsg
}

//helping function for msgHandler to check if all IDs in peers have confirmed
func leftSubsetOfRight(left []string, right []string) bool {
	isSubset := true
	for i := 0; i < len(left); i++ {
		containsElement := false
		for j := 0; j < len(right); j++ {
			if right[j] == left[i] {
				containsElement = true
			}
		}
		if !containsElement {
			isSubset = false
		}
	}
	return isSubset
}

//single handshake message handler with peers
func txMsgHandler(txCh chan types.Message, msgCh chan types.Message, thisElevID string,
	receivedUpdateCh chan types.Message) {
	//elevator will always be among peers
	for {
		select {
		case newMsg := <-msgCh:
			//set received confirmations to only this elevator
			confirmationMutex.Lock()
			peerMutex.Lock()
			currentConfirmations = append(currentConfirmations, thisElevID)
			currentConfirmations = currentConfirmations[(len(currentConfirmations) - 1):]
			//resend until confirmed by everyone listed on network
			for !leftSubsetOfRight(currentPeers, currentConfirmations) {
				confirmationMutex.Unlock()
				peerMutex.Unlock()
				txCh <- newMsg
				time.Sleep(15 * time.Millisecond)
				confirmationMutex.Lock()
				peerMutex.Lock()
			}
			confirmationMutex.Unlock()
			peerMutex.Unlock()
			receivedUpdateCh <- newMsg
		}
	}
}

//handles incoming messages.
func rxMsgHandler(rxCh chan types.Message, peerChannel chan peers.PeerUpdate, thisElevID string,
	peerDistributorCh chan peers.PeerUpdate, toDistributorCh chan []types.Elevator, txCh chan types.Message) {
	var lastReceivedID map[string]int
	lastReceivedID = make(map[string]int)
	lastReceivedID[thisElevID] = 1000000000000000000
	for {
		select {
		case peerUpdate := <-peerChannel:
			if len(peerUpdate.New) > 0 {
				lastReceivedID[peerUpdate.New] = 0
			}
			if len(peerUpdate.Peers) == 0 {
				peerMutex.Lock()
				currentPeers = append(currentPeers, thisElevID)
				peerMutex.Unlock()
			}
			if len(peerUpdate.Peers) > 0 {
				peerMutex.Lock()
				currentPeers = peerUpdate.Peers
				peerMutex.Unlock()
			}
			peerDistributorCh <- peerUpdate
		case receivedMsg := <-rxCh:

			if receivedMsg.ConfirmedMsgOwner == thisElevID {
				confirmationMutex.Lock()
				currentConfirmations = append(currentConfirmations, receivedMsg.ElevID)
				confirmationMutex.Unlock()
			}
		//	fmt.Printf("received msg id %v \n", receivedMsg.MsgID)
		//	fmt.Printf("last received id %v \n", lastReceivedID[receivedMsg.ElevID] )
			if receivedMsg.MsgID > lastReceivedID[receivedMsg.ElevID] {
				toDistributorCh <- receivedMsg.Content
				sendConfirmation(receivedMsg.ElevID, thisElevID, txCh)
			} else if receivedMsg.MsgID != 0 {
				sendConfirmation(receivedMsg.ElevID, thisElevID, txCh)
			}
		}
	}
}

//Network is the function to be called from outside the module
func Network(fromDistributorCh chan []types.Elevator, 
			toDistributorCh chan []types.Elevator,
			peerTxEnableCh chan bool, 
			peerDistributorCh chan peers.PeerUpdate, 
			myID chan string,
			portNum string,
			ID string) {
				
	peerElevPort := 2203
	elevMsgPort := 1903
	//initialize channels
	peerCh := make(chan peers.PeerUpdate)
	txCh := make(chan types.Message)
	rxCh := make(chan types.Message)
	msgCh := make(chan types.Message, 10)
	receivedUpdateCh := make(chan types.Message)
	msgID := 1
	// now := time.Now()int
	// ID := now.UnixNano()lightCh
	// localIP := strconv.FormatInt(ID, 10)
	// myID <- localIP
	
	
	//localIP := portNum
	localIP := ID
	// if err != nil {
	// 	fmt.Println(err)
	// 	localIP = "DISCONNECTED"
	// }
	myID <-localIP


	//start go routines
	go peers.Transmitter(peerElevPort, localIP, peerTxEnableCh)
	go peers.Receiver(peerElevPort, peerCh)
	go bcast.Transmitter(elevMsgPort, txCh)
	go bcast.Receiver(elevMsgPort, rxCh)
	go txMsgHandler(txCh, msgCh, localIP, receivedUpdateCh)
	go rxMsgHandler(rxCh, peerCh, localIP, peerDistributorCh, toDistributorCh, txCh)
	for {
		select {
		//new message from dist
		case msg := <-fromDistributorCh:
			newMsg := composeMsg(localIP, msg, msgID)
			msgCh <- newMsg
			msgID++
		//new update from rx
		case newUpdate := <-receivedUpdateCh:
			//fmt.Printf("Sending forom network \n")
			toDistributorCh <- newUpdate.Content
		}
	}
}
