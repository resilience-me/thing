# Ripple + Resilience, in a very simple true peer-to-peer implementation

Uses UDP, single-threaded server, a standardized datagram format for both client and server (and a datagram for server-to-client responses). Any data stored in alphanumeric format in text files. Accounts are identified by a username, and, the address of their host server (IP address or domain name). "Database" managed with simple directories, `datadir/resilience/accounts/username/peers/server_address/username`. No encryption, only authentication. Counters (nonces) to prevent replay attacks. Retransmission and verification that payment succeeded and such can be added manually via client, this repository is the server only.

    type Datagram struct {
        Command        byte
        XUsername      [32]byte
        YUsername      [32]byte
        YServerAddress [32]byte
        Arguments      [256]byte
        Counter        [4]byte
        Signature      [32]byte
    }

The command is one byte, allowing 256 commands. The first 128 commands are client commands, the last 128 are server commands. In datagram, user X and user Y mean opposite things when using client commands and server commands. In client commands, the Y user is the user you want to interact with, and the X user is yourself. In server commands, the Y user is the user whose account sent the command, and the X user is the target user that the account wants to interact with. The counter is managed with different counters for different commands. The signature relies on a symmetric secret key, in client command shared by the server and the client, and in server commands shared by two users with a direct connection in the system (they will be user X and user Y in a given datagram. ) It is a simple sha256 of the datagram together with the secret key. And, the 256 byte long arguments field, can hold arbitrary data for operands to the command. The datagram is 389 bytes and fits well within a UDP packet.

Client return values use a standardized datagram too. It includes a random "nonce" (so client can decide to listen to specific response), and then an a byte array Result with a return value or some arbitrary return data. The nonce is the signature of the datagram sent by the client. The response datagram also includes a signature.

    type ResponseDatagram {
        Nonce        [32]byte
        Result       [256]byte
        Signature    [32]byte
    }

### Handling trustlines

A number of counters keep track of state of trustlines. There is `counter`, that tracks an account's most recent update to their trustline to another acount. There is `sync_in` and `sync_out`, that track synchronization of trustlines. There is also `timestamp`, for an account to locally track when an incoming trustline was last synced. Then, to manage updating the sync timestamp even if there was no need to sync (but an attempt was made), there is `sync_counter_in` and `sync_counter_out`. The platform does not use timestamps as counters or "nonces", but to allow tracking most recent change locally, the extra sync counters are added.
