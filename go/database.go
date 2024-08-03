import (
    "os"
    "path/filepath"
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// GetPeerDir constructs the peer directory path from the datagram.
func GetPeerDir(dg Datagram) string {
    username := string(dg.XUsername[:])
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    return filepath.Join(datadir, "accounts", username, "peers", peerAddress, peerUsername)
}
