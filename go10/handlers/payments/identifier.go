package payments

import (
	"fmt"
	"crypto/sha256"
	"ripple/types"
	"ripple/pathfinding"
	"ripple/config"
)

func concatNameAndServer(username, serverAddress string) []byte {
  return append(types.PadStringTo32Bytes(username), types.PadStringTo32Bytes(serverAddress)...)
}

func generatePaymentIdentifier(dg *types.Datagram, inOrOut byte) string {
  user := concatNameAndServer(dg.Username, config.GetServerAddress())
  peer := concatNameAndServer(dg.PeerUsername, dg.PeerServerAddress)
  
  var preimage []byte
  
  if inOrOut == types.Incoming {
    preimage = append(peer, user...)
  } else {
    preimage = append(user, peer...)
  }
  preimage = append(preimage, dg.Arguments[:8]...)
  hash := sha256.Sum256(preimage)
  
  return fmt.Sprintf("%x", hash[:])
}


// GenerateAndInitiatePayment handles the generation of the payment identifier and initiation of the payment.
func GenerateAndInitiatePayment(datagram *types.Datagram, inOrOut byte) {
    // Generate the Payment struct for an incoming payment
    identifier := generatePaymentIdentifier(datagram, inOrOut)
    nonce := types.BytesToUint32(datagram.Arguments[4:8])
    payment := pathfinding.NewPayment(datagram, identifier, inOrOut, nonce)
    amount := types.BytesToUint32(datagram.Arguments[0:4])
    // Initiate the incoming payment using the constructed Payment struct
    pathfinding.GetPathManager().InitiatePayment(datagram.Username, payment, amount)
}
