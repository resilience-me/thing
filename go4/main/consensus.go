// HandleTransactionRequest processes a transaction request from the non-validator.
func HandleTransactionRequest(filename string, request []byte, validatorID []byte) error {
    // Ensure that this account is the current validator
    ValidatedLatestBlock, err := IsValidator(filename, validatorID)
    if err != nil {
        return err
    }
    if ValidatedLatestBlock {
        return fmt.Errorf("this account is not the current validator")
    }

    // Convert the request into a full transaction
    rawTransaction, err := ConvertRawBytesToTransaction(request)
    if err != nil {
        return fmt.Errorf("failed to convert request to transaction: %v", err)
    }

    // Convert the request into a full transaction
    rawTransaction, err := PrepareTransaction(filename, request, validatorID)
    if err != nil {
        return fmt.Errorf("failed to prepare transaction: %v", err)
    }

    // Validate and append the transaction
    err = SubmitTransaction(filename, rawTransaction, validatorID)
    if err != nil {
        return fmt.Errorf("failed to validate and add transaction: %v", err)
    }

    return nil
}

// ConvertRawBytesToTransaction converts raw bytes of a TransactionRequest to a Transaction by populating the fields.
func ConvertRawBytesToTransaction(request []byte) ([]byte, error) {
    // Create a byte slice to hold the transaction
    rawTransaction := make([]byte, LengthTransaction)

    // Copy the entire TransactionRequest data into the rawTransaction starting at OffsetFrom
    copy(rawTransaction[OffsetFrom:], request[:SizeRequest-SizeSignature])

    return rawTransaction, nil
}
