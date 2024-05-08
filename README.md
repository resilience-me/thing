# Ryan Fugger's Ripple in P2P way

Since Ripple, uniquely, can rely only on authentication between the two people making the exchange, this implementation uses authentication only between people, and not between servers.

What this means is, instead of each server authenticating itself via a certificate authority, which requires a centralization of trust (even if a web-of-trust certificate authority is used, it still to some extent centralizes), Ripple can operate with only people with trustlines authenticating themselves to each other.

Practically, this also works well with symmetric keys, which the people exchange in some way. A peer (essentially a trustline, but with some other data too such as shared secret key) is stored in `accounts/your_account/peers/your_peers_account/secret_key.txt`, and the shared secret is stored there as well. The authentication then uses a "message authentication code" alongside the message, a hash. Such a signature has theoretically stronger security than an asymmetric signature, and this implementation wants to demonstrate that Ripple can operate with symmetric cryptography only, which is a strength it has.

This implementation will then use no encryption of the messages. It will use no encryption since encryption isn't strictly needed to run Ripple. It is easy to add. Note that assuming account-to-account encryption, the account identifier has to be plaintext anyway.

People also use symmetric authentication with their server, and this is set up by exchanging a shared secret key with the server admin. The key is stored (on the server) in `accounts/your_account/secret_key.txt`, and in the client, in `client_datadir/secret_key.txt`. Besides that, all messages in plaintext. No persistent connection to server needed, craft a message (a command with argments, and your username as parameter), generate hash as signature, and submit the message and the signature to the server. Asymmetric key could be used too, but the benefit of asymmetric cryptography is in public contexts, and in person-to-person (including person-to-server where its still a personal exchange between two entities) they're not required.

The system can probably run over UDP, and be based on broadcast, and if the frame was not delivered, the ability to poll for if the command was processed. All commands may fit within a single frame, making it very simple. A tentative format for a datagram in the system:

    typedef struct {
        uint8_t connectionType;  // Type of connection: 0 for client, 1 for server, etc.
        union {
            ClientDatagram client;
            ServerDatagram server;
        } data;                   // Union holding either client or server datagram
    } Datagram;
    
    typedef struct {
        char sourceUser[32];      // Username, up to 31 characters + null terminator
        uint8_t command;          // Numeric code representing the command
        char arguments[256];      // Arguments for the command
        char signature[32];       // SHA-256 hash signature for verification
    } ClientDatagram;
    
    typedef struct {
        char sourceDomain[32];    // Domain of the source server
        char sourceUser[32];      // Source account username
        char destinationUser[32]; // Target account username
        uint8_t command;          // Operation code
        char arguments[256];      // Command arguments
        char signature[32];       // Signature for verification
    } ServerDatagram;

### Commands

Tentatively, the opcodes will be divided between client and server opcodes. 0-127 will be client opcodes and 128-255 will be server opcodes.

    1. SET_TRUSTLINE
    Value: 0x01
    Description: Sets or updates a trustline to a person.
    Parameters Encoding:
    username (32 byte)
    size (64 byte)
    
    2. GET_TRUSTLINE
    Value: 0x02
    Description: Retrieves size of trustline to a person.
    Parameters Encoding:
    username (32 byte)
    
    [...]
    
    128. SET_TRUSTLINE
    Value: 0x80
    Description: Synchronize trustline update between two accounts
    Parameters Encoding:
    username (32 byte)
    size (64 byte)
