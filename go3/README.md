# Ripple in a very simple true peer-to-peer implementation

Uses TCP to transfer single-packet Datagrams with "guaranteed delivery" (the retransmission built-into TCP. ) Communication is one-way in server-to-server exchanges (thus a single Datagram transferred, then disconnect), and two-way in client-to-server (allows the client to receive a response such as error or success. ) An account processes one Datagram at a time (coordinated via SessionManager class. ) Any data stored in alphanumeric format in text files. Accounts are identified by a username, and, the address of their host server (IP address or domain name). "Database" managed with simple directories, `datadir/ripple/accounts/username/peers/server_address/username`. No encryption, only authentication. This repository is the server only.

    type Datagram struct {
        Command           byte
        Username          [32]byte
        PeerUsername      [32]byte
        PeerServerAddress [32]byte
        Arguments         [256]byte
        Counter           [4]byte
        Signature         [32]byte
    }

The command is one byte, allowing 256 commands. The first 128 commands are client commands, the last 128 are server commands. The counter is managed with different counters for different commands. The signature relies on a symmetric secret key, in client command shared by the server and the client, and in server commands shared by two users with a direct connection in the system. It uses HMAC. And, the 256 byte long arguments field, can hold arbitrary data for operands to the command. The datagram is 389 bytes.

### Handling trustlines

A number of counters keep track of state of trustlines. There is `counter`, that tracks an account's most recent update to their trustline to another acount. There is `sync_in` and `sync_out`, that track synchronization of trustlines (relative to `counter`). There is also `timestamp`, for an account to locally track when an incoming trustline was last synced. Then, to manage updating the sync timestamp even if there was no need to sync (but an attempt was made), there is `sync_counter_in` and `sync_counter_out`. The platform does not use timestamps as counters or "nonces", but to allow tracking most recent change locally, the extra sync counters are added.

### Coordinating payments

Step 1) A command to place a time lock on the trustlines is sent down the path. 

Step 2) A command to finalize the commit is sent down the path. This increases the time lock, and, adds a rule that the commit can only be aborted if it is verified that the next in line has aborted it, or never received it. (Thus if it reaches buyer, it cannot be cancelled unless buyer somehow cancels it... )

Step 3) A command to finalize the payment is sent down the path. A credit line has now formed, and the payment is complete.
