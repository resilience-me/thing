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
