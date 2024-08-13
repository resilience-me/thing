// GetPeers retrieves a list of all peer accounts for a given username
func GetPeers(username string) ([]PeerAccount, error) {
    var peers []PeerAccount
    baseDir := filepath.Join("datadir", "accounts", username, "peers")

    // Read all server address directories in the peers directory
    serverDirs, err := ioutil.ReadDir(baseDir)
    if err != nil {
        return nil, fmt.Errorf("unable to read directory %s: %v", baseDir, err)
    }

    // Iterate over all server address directories
    for _, serverDir := range serverDirs {
        if serverDir.IsDir() {
            serverAddress := serverDir.Name()
            serverPath := filepath.Join(baseDir, serverAddress)

            // Read all peer directories under the current server address
            peerDirs, err := ioutil.ReadDir(serverPath)
            if err != nil {
                return nil, fmt.Errorf("unable to read directory %s: %v", serverPath, err)
            }

            // Iterate over all peer directories and create PeerAccount structs
            for _, peerDir := range peerDirs {
                if peerDir.IsDir() {
                    peer := PeerAccount{
                        Username:      peerDir.Name(),
                        ServerAddress: serverAddress,
                    }
                    peers = append(peers, peer)
                }
            }
        }
    }

    return peers, nil
}