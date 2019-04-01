package communication

import (
	"sync"
	"time"
	"../peers"
	"../../types"
)

var currentConfirmations []string
var currentPeers []string
var confirmationMutex = &sync.Mutex{}
var peerMutex = &sync.Mutex{}

func ComposeMsg(ID string, content []types.Elevator, msgID int) types.Message {
	var newMsg types.Message
	newMsg.ElevID = ID
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
func TxMsgHandler(txCh chan types.Message, msgCh chan types.Message, thisElevID string,
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
func RxMsgHandler(rxCh chan types.Message, peerChannel chan peers.PeerUpdate, thisElevID string,
	peerDistributorCh chan peers.PeerUpdate, toDistributorCh chan []types.Elevator, txCh chan types.Message) {
	var lastReceivedID map[string]int
	lastReceivedID = make(map[string]int)
	lastReceivedID[thisElevID] = 100000000000
	for {
		select {
		case peerUpdate := <-peerChannel:
			if (len(peerUpdate.New) > 0) && (peerUpdate.New != thisElevID) {
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
			if receivedMsg.MsgID > lastReceivedID[receivedMsg.ElevID] {
				toDistributorCh <- receivedMsg.Content
				sendConfirmation(receivedMsg.ElevID, thisElevID, txCh)
			} else if receivedMsg.MsgID != 0 {
				sendConfirmation(receivedMsg.ElevID, thisElevID, txCh)
			}
		}
	}
}

