import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
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

func hashAndTruncateToAddress() []byte {
    hash := sha256.Sum256(combined)
    return hash[:20]  // Use the first 20 bytes of the hash
}
// GenerateAddress generates an address using SHA-256 from an ECDSA public key, skipping the prefix.
func GenerateAddress(pubKey *ecdsa.PublicKey) []byte {
    // Get the uncompressed public key bytes
    pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
    // The address is the first 20 bytes of the hash
    return hashAndTruncateToAddress(pubKeyBytes[1:])
}
// generateIdentifier creates a unique identifier using the first 20 bytes of a SHA-256 hash of the combined addresses.
func generateIdentifier(myAddress, otherAddress []byte) []byte {
    combined := append(myAddress, otherAddress...)
    return hashAndTruncateToAddress(combined)
}
