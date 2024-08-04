# Ripple + Resilience, in a very simple true peer-to-peer implementation

Uses UDP, single-threaded server, a standardized datagram format for both client and server (and a datagram for server-to-client responses). Any data stored in alphanumeric format in text files. Accounts are identified by a username, and, the address of their host server (IP address or domain name). "Database" managed with simple directories, `datadir/resilience/accounts/username/peers/server_address/username`. No encryption, only authentication. The reason to skip encryption is because it is not indispensable for the system. And the threat of eavesdropping is quite low (see ["On skipping encryption"](https://github.com/resilience-me/thing/blob/main/go/README.MD#on-skipping-encryption) below. ) Counters (nonces) to prevent replay attacks. Retransmission and verification that payment succeeded and such can be added manually via client, this repository is the server only. Client can be "extended" with plug-in asymmeteric client to allow simple login via something like a website with password.

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

Client return values use a standardized datagram too. It includes a random "nonce" (so client can decide to listen to specific response), and then an a byte array Result with a return value or some arbitrary return data. The response datagram also includes a signature.

    type ResponseDatagram {
        Nonce        [4]byte
        Result       [256]byte
        Signature    [32]byte
    }

### On skipping encryption

Account messages are not that sensitive information. Many people would be OK with having a conversation in public about their trustline to a friend, and eavesdropping on their internet routers and communication lines is comparable to eavesdropping in public places like cafés, a town square or a park. A "message authentication code" such as a hash function is also theoretically more secure than encryption (it is not reversible so theoretically more secure. ) It is easy to add encryption, either via asymmetric encryption, or, something like a one-time pad (the latter gives perfect secrecy, a property that Ripple is uniquely able to have. ) It is skipped in this implementation to keep the architecture simple, and because many may not actually care that much about eavesdropping attack vectors (eavesdropping is not infinitely simple anyway. ) And there's also often encryption at "link layer" such as in mobile networks, and there's probably transport protocols with encryption that could easily substitute the UDP here (rather than having to add encryption in the "application layer.. )
