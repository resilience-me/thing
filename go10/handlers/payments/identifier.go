package payments

import (
	"fmt"
	"crypto/sha256"
	"ripple/types"
)

func concatAccounts(username, serverAddress string) []byte {
  return append(types.PadStringTo32Bytes(), types.PadStringTo32Bytes())
}

func generatePaymentIdentifier(dg *Datagram, inOrOut byte) string {
  user := concatAccounts(dg.Username, GetServerAddress())
  peer := concatAccounts(dg.PeerUsername, dg.PeerServerAddress)
  
  var preimage []byte
  
  if inOrOut == types.Incoming {
    preimage = append(peer, user)
  } else {
    preimage = append(user, peer)
  }
  preimage = append(preimage, dg.Arguments[:8])
  hash := sha256.Sum256(preimage)
  
  return fmt.Sprintf("%x", hash[:])
}


// GenerateAndInitiatePayment handles the generation of the payment identifier and initiation of the payment.
func GenerateAndInitiatePayment(session main.Session, inOrOut byte) {
    // Generate the Payment struct for an incoming payment
    identifier := generatePaymentIdentifier(session.Datagram, inOrOut)
    payment := pathfinding.NewPayment(datagram, identifier, inOrOut)

    // Initiate the incoming payment using the constructed Payment struct
    session.PathManager.InitiatePayment(session.Datagram.Username, payment)
}
