import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// GenerateAddress generates an address using SHA-256 from an ECDSA public key, skipping the prefix.
func GenerateAddress(pubKey *ecdsa.PublicKey) string {
    // Get the uncompressed public key bytes
    pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

    // Compute the SHA-256 hash of the public key bytes, skipping the 0x04 prefix
    hash := sha256.Sum256(pubKeyBytes[1:]) // Skip the 0x04 prefix

    // The address is the last 20 bytes of the hash
    address := hash[len(hash)-20:]

    // Return the address as a hexadecimal string prefixed with 0x
    return fmt.Sprintf("0x%x", address)
}
