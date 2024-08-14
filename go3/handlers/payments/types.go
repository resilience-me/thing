type PaymentDetails struct {
    Counterpart PeerAccount
    InOrOut     bool
    Amount      uint32
    Nonce       uint32
}
