# Ripple in a very simple true peer-to-peer implementation

Custom transport protocol, UDP + retransmission and acknowledgement over ephemeral port. At the application layer, counters to prevent datagrams from being replayed. No encryption, only authentication. An account processes one Datagram at a time (coordinated via SessionManager class. ) Accounts are identified by a username, and, the address of their host server (IP address or domain name). "Database" managed with simple directories, `datadir/accounts/username/peers/server_address/username`. Any data stored in alphanumeric format in text files. This repository is the server only.

    type Datagram struct {
        Command           byte
        Username          [32]byte
        PeerUsername      [32]byte
        PeerServerAddress [32]byte
        Arguments         [256]byte
        Counter           [4]byte
        Signature         [32]byte
    }

The command is one byte, allowing 256 commands. The first 128 commands are client commands, the last 128 are server commands. The signature relies on a symmetric secret key, in client command shared by the server and the client, and in server commands shared by two users with a direct connection in the system. It uses sha256. And, the 256 byte long arguments field can hold arbitrary data for operands to the command. The datagram is 389 bytes.

### Counters

There is three main sets of counters to prevent datagrams being replayed. One for client to server interactions (`counter.txt` in `accounts/username`), and two for server to server interactions (one per direction) for each peer account a user account has (`counter_out.txt` and `counter_in.txt` in `accounts/username/peers/server_address/username`).

### Handling trustlines

A number of counters keep track of state of trustlines. There is "sync counter", that tracks how many times the trustline has been updated. And, `sync_in` and `sync_out`, that track synchronization of trustlines (relative to `sync_counter`). There is also `timestamp`, for an account to locally track when an incoming trustline was last synced. The timestamp is never exchanged and there is no need for consensus on time, the platform does not use timestamps as counters or "nonces".

### Path finding

The Path finding is very simple. It is practically “stateless”, no routing tables are stored, all routing is generated for each payment request.

The path-finding optimizes for never going too deep. It is bidirectional, reducing accounts queried to 2*sqrt(unidirectional). And, it searches in increments of 1, always returning to the root before increasing the depth by 1. Thus, whenever a path is found, the search ends (the root stops replying to response by incrementing request. ) Path requests use an identifier that is a simple random number, and are sent both from buyer and receiver. Whenever these “fronts” meet, a path is found, and the first path found is chosen. The "first path found" approximates fewest hops.

### Coordinating payments

Step 1) A command to place a time lock on the trustlines is sent down the path. 

Step 2) A command to finalize the commit is sent down the path. This increases the time lock, and, adds a rule that the commit can only be aborted if it is verified that the next in line has aborted it, or never received it. (Thus if it reaches buyer, it cannot be cancelled unless buyer somehow cancels it... )

Step 3) A command to finalize the payment is sent down the path. A credit line has now formed, and the payment is complete.
