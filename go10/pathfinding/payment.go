package pathfinding

func (pm *PathManager) CleanupCacheAndFetchAccount(username string) *Account {
    // Cleanup all accounts first
    pm.Cleanup()

    account := FetchAndRefresh()
    if != nil {
        account.Cleanup()
        return account
    }
    return pm.Add(username)
}

// InitiatePayment sets up or updates payment details for an account, creating the account if necessary.
func (pm *PathManager) InitiatePayment(username string, payment *Payment, amount uint32) {
    // Fetch or create the account, with any necessary cleanup
    account := pm.CleanupCacheAndFetchAccount(username)

    // If a previous payment existed, remove it
    if account.Payment != nil {
        account.Remove(account.Payment.Identifier)
    }

    // Set or update the payment details
    account.Payment = payment

    // Add or update the related Path entry with a new timestamp
    account.Add(payment.Identifier, amount, PeerAccount{}, PeerAccount{})  // No PeerAccount details needed
}
