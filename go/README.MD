# Ripple + Resilience, in a very simple true peer-to-peer implementation

Uses UDP, single-threaded server, a standardized datagram format for both client and server. Any data stored in alphanumeric format in text files. Accounts are identified by a username, and, the address of their host server (IP address or domain name). "Database" managed with simple directories, `datadir/resilience/accounts/username/peers/server_address/username`. No encryption, only authentication (as eavesdropping can happen in daily life too. ) Client can be "hidden" behind plug-in asymmeteric client by using `LocalClientMode`, to prevent anyone from accessing the getter functions. Otherwise all data is public to anyone who bothers to read from it. Counters (nonces) to prevent replay attacks. Retransmission and verification that payment succeeded and such can be added manually via client, this repository is the server only.

    type Datagram struct {
        Command        byte
        XUsername      [32]byte
        YUsername      [32]byte
        YServerAddress [32]byte
        Arguments      [256]byte
        Counter        [4]byte
        Signature      [32]byte
    }

The command is one byte, allowing 256 commands. The first 128 commands are client commands, the last 128 are server commands. In datagram, user X and user Y mean opposite things when using client commands and server commands. In client commands, the Y user is the user you want to interact with, and the X user is yourself. In server commands, the Y user is the user whose account sent the command, and the X user is the target user that the account wants to interact with. The counter is managed with different counters for different commands. The signature relies on a symmetric secret key, in client command, shared by the server and the client, and in server commands shared by two users with a direct connection in the system (that are user X and user Y. )

is a simple sha256 of the datagram together with 
