import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/sha256"
)

const (
    SizeNumber      = 4
    SizeValidator   = 20
    SizeFrom        = 20
    SizeTo          = 20
    SizeData        = 256
    SizeParentHash  = 32
    SizeSignature   = 64

    SizeTransaction = 416
    SizeRequest     = 360

    OffsetNumber    = 0
    OffsetValidator = OffsetNumber + SizeNumber
    OffsetFrom      = OffsetValidator + SizeValidator
    OffsetTo        = OffsetFrom + SizeFrom
    OffsetData      = OffsetTo + SizeTo
    OffsetParentHash = OffsetData + SizeData
    OffsetSignature = OffsetParentHash + SizeParentHash
)

type Transaction struct {
    Number            [4]byte
    Validator         [20]byte
    From              [20]byte
    To	              [20]byte
    Data              [256]byte
    ParentHash        [32]byte
    Signature         [64]byte
}

type TransactionRequest struct {
    From      [20]byte
    To        [20]byte
    Data      [256]byte
    Signature [64]byte
}

type Datagram struct {
    Identifier  [20]byte
    Ciphertext	[360]byte
}

// GenerateAddress generates an address using SHA-256 from an ECDSA public key, skipping the prefix.
func GenerateAddress(pubKey *ecdsa.PublicKey) []byte {
    // Get the uncompressed public key bytes
    pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
    // The address is the first 20 bytes of the hash
    hash := sha256.Sum256(pubKeyBytes[1:])
    return hash[:20]  // Use the first 20 bytes of the hash
}
// GenerateIdentifier creates a unique identifier using the first 20 bytes of a SHA-256 hash of the combined addresses.
func GenerateIdentifier(myAddress, otherAddress []byte) []byte {
    return XORBytes(myAddress, otherAddress)
}

// XORBytes takes two byte slices and returns their XOR combination.
func XORBytes(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}
