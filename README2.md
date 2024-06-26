A very simple or "bare minimum" implementation of Ryan Fugger's Ripple. Client-server architecture, where users host accounts on servers and interact with their server via a client, and accounts on servers can interact both within and between servers. Single-threaded, uses UDP. A form of "retransmission" possible by manual triggering from the user client, used only for committing and finalizing payments. Message authentication using a sha-256 hash. Users agree on a shared secret with their server manually (how is up to the user and server), and use it for message authentication. When accounts interact between servers, they also rely on a shared secret to generate sha256 message authentication codes. Thus account pairs also need shared secrets, and these are generated by the peers manually, and uploaded to their servers. When uploading a peer secret key, encryption is used, using the client-server shared secret. The encryption is with a derived key, sha256 hash of the secret key and a nonce (incremented from 0 to N), and then a simple XOR operation. Since encryption is only strictly used when submitting peer secret keys, it is the only place where encryption is used. "Database" is by storing data in files, the file tree is datadir in `~/.ripple/`, client data in `datadir/client` and server data in `datadir/server`. In server folder accounts are in `accounts/username`, and in an account folder peers are stored in `peers/server_address/username`. Thus secret key for client-user is in `secretkey.txt` in the client directory, and in `accounts/username/secretkey.txt` in the server directory, and within account directories the secret keys for a peer is in `peers/server_address/username/secretkey.txt`.
