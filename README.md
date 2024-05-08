# Ryan Fugger's Ripple in P2P way

Since Ripple, uniquely, can rely only on authentication between the two people making the exchange, this implementation uses authentication only between people, and not between servers.

What this means is, instead of each server authenticating itself via a certificate authority, which requires a centralization of trust (even if a web-of-trust certificate authority is used, it still to some extent centralizes), Ripple can operate with only people with trustlines authenticating themselves to each other.

Practically, this also works well with symmetric keys, which the people exchange in some way. A peer (essentially a trustline, but with some other data too such as shared secret key) is stored in `accounts/your_account/peers/your_peers_account/secret_key.txt`, and the shared secret is stored there as well. The authentication then uses a "message authentication code" alongside the message, a hash. Such a signature has theoretically stronger security than an asymmetric signature, and this implementation wants to demonstrate that Ripple can operate with symmetric cryptography only, which is a strength it has.

This implementation will then use no encryption of the messages. It will use no encryption since encryption isn't strictly needed to run Ripple. It is easy to add. Note that assuming account-to-account encryption, the account identifier has to be plaintext anyway.

People also use symmetric authentication with their server, and this is set up by exchanging a shared secret key with the server admin. The key is stored (on the server) in `accounts/your_account/secret_key.txt`, and in the client, in `client_datadir/secret_key.txt`. Besides that, all messages in plaintext. No persistent connection to server needed, craft a message (a command with argments, and your username as parameter), generate hash as signature, and submit the message and the signature to the server. Asymmetric key could be used too, but the benefit of asymmetric cryptography is in public contexts, and in person-to-person (including person-to-server where its still a personal exchange between two entities) they're not required.

The system can probably run over UDP, and be based on broadcast, and if the frame was not delivered, the ability to poll for if the command was processed. All commands may fit within a single frame, making it very simple. A tentative format for a datagram in the system:

    typedef struct {
        uint8_t connectionType;    // 0 for client-server, 1 for server-server interactions.
        char x_username[32];       // Username for user X, context-dependent.
        char y_username[32];       // Username for user Y, context-dependent.
        char y_domain[32];         // Domain of user Y, context-dependent.
        uint8_t command;           // Numeric code for the command to be executed.
        char arguments[256];       // Data necessary for executing the command.
        uint32_t nonce;            // 4-byte nonce field for replay protection, context-dependent.
        char signature[32];        // SHA-256 hash for verifying data integrity and authenticity.
    } Datagram;

In most client-to-server interactions as well as server-to-server interactions, two users are involved. One of the users is on the server that receives the datagram, and thus organized under "localhost", whereas the other needs a domain name as part of their identifier. These are "user X" and "user Y" in the datagram, where "user Y" also has a domain name identifier. When a user interacts with a server via a client, user X will be their account, and user Y will be the account they may want to interact with (such as setting the trustline for. ) And vice versa, when a server interacts with another server (on behalf of a user account), "user Y" will be their account and they also provide a domain name as part of the identifier, and user X will be the account they want to interact with.

Domain name can of course be fetched via reverse DNS lookup, but it seems simpler to just pass it with the datagram, since it is part of the user idenfifier information.

The nonce is either between user (client) and server, or per account relationship in server-to-server.

### Commands

Client commands:

    1. SET_TRUSTLINE
    Value: 0x01
    Description: Sets or updates a trustline to a person.
    Arguments:
    size (64 byte)
    
    2. GET_TRUSTLINE
    Value: 0x02
    Description: Retrieves size of trustline to a person.
    Arguments Encoding:

Server commands:
    
    1. SET_TRUSTLINE
    Value: 0x01
    Description: Synchronize trustline update between two accounts
    Arguments Encoding:
    size (64 byte)
