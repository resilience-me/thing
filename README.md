# Ryan Fugger's Ripple in P2P way

Since Ripple, uniquely, can rely only on authentication between the two people making the exchange, this implementation uses authentication only between people, and not between servers.

What this means is, instead of each server authenticating itself via a certificate authority, which requires a centralization of trust (even if a web-of-trust certificate authority is used, it still to some extent centralizes), Ripple can operate with only people with trustlines authenticating themselves to each other.

Practically, this also works well with symmetric keys, which the people exchange in some way. A peer (essentially a trustline, but with some other data too such as shared secret key) is stored in `accounts/your_account/peers/your_peers_account/secretkey.txt`, and the shared secret is stored there as well. The authentication then uses a "message authentication code" alongside the message, a hash. Such a signature has theoretically stronger security than an asymmetric signature, and this implementation wants to demonstrate that Ripple can operate with symmetric cryptography only, which is a strength it has.

This implementation will then use no encryption of the messages. It will use no encryption since encryption isn't strictly needed to run Ripple. It is easy to add. Note that assuming account-to-account encryption, the account identifier has to be plaintext anyway.

People also use symmetric authentication with their server, and this is set up by exchanging a shared secret key with the server admin. The key is stored (on the server) in `accounts/your_account/secretkey.txt`, and in the client, in `client_datadir/secretkey.txt`. Besides that, all messages in plaintext. No persistent connection to server needed, craft a message (a command with argments, and your username as parameter), generate hash as signature, and submit the message and the signature to the server. Asymmetric key could be used too, but the benefit of asymmetric cryptography is in public contexts, and in person-to-person (including person-to-server where its still a personal exchange between two entities) they're not required.

The system can probably run over UDP, and be based on broadcast, and if the frame was not delivered, the ability to poll for if the command was processed. All commands may fit within a single frame, making it very simple. A form of retransmission can be done "manually" from the client end, with somelike like an RETRY_PAYMENT command to retransmit if payment got stuck.

A tentative format for a datagram in the system:

    typedef struct {
        uint8_t command;           // Command code with most significant bit specifying connection type (client/server)
        char x_username[32];       // Username for user X, context-dependent.
        char y_username[32];       // Username for user Y, context-dependent.
        char y_server_address[32]; // Server address of user Y, context-dependent.
        char arguments[256];       // Data necessary for executing the command.
        uint32_t nonce;            // 4-byte nonce field for replay protection, context-dependent.
        char signature[32];        // SHA-256 hash for verifying data integrity and authenticity.
    } Datagram;

The connection type and command is used to dispatch the command handler for a command:

    typedef void (*CommandHandler)(const Datagram*, int sockfd, const struct sockaddr_storage*);
    CommandHandler command_handlers[256] = {
        [0] = client_handle_set_trustline,
        [1] = client_handle_get_trustline,
        [128] = server_handle_set_trustline
    };


The command handlers are dispatched as `command_handlers[dg->command]`. Command codes 0 to 127 are for client commands and 128 to 255 are for server commands (thus the most significant bit specifies connection type there, and this bit is also used when verifying the signature since user x and user y parameters have to be used differently depending on client or server interaction context. ) The command handlers manage potential response to client (or server).

To get the connection type, use a simple bitshift:

    int getConnectionType(const Datagram *dg) {
        return (dg->command >> 7) & 1;
    }

And to see the command handler dispatcher in the context of the while loop that accepts UDP datagrams:

    while (1) {
        // Receive data
        ssize_t recv_len = recvfrom(sockfd, buffer, BUFFER_SIZE, 0, (struct sockaddr*)&client_addr, &addr_len);
        if (recv_len == -1) {
            continue;
        }

        // Deserialize received data
        Datagram dg;
        deserialize_datagram(buffer, &dg);

        // Verify signature and nonce
        if (!verify_signature(buffer, &dg) || !verify_nonce(dg.nonce)) continue;

        // Call appropriate command handler
        CommandHandler handler = command_handlers[dg.command];
        if (handler) {
            handler(&dg, sockfd, *(struct sockaddr_in *)&client_addr);
        }
    }

In most client-to-server interactions as well as server-to-server interactions, two users are involved. One of the users is on the server that receives the datagram, and thus organized under "localhost", whereas the other needs a server address as part of their identifier. These are "user X" and "user Y" in the datagram, where "user Y" also has a server address identifier. When a user interacts with a server via a client, user X will be their account, and user Y will be the account they may want to interact with (such as setting the trustline for. ) And vice versa, when a server interacts with another server (on behalf of a user account), "user Y" will be their account and they also provide a server address name as part of the identifier, and user X will be the account they want to interact with.

The server address can of course be read from the IP header and if it was a domain it can be fetched via reverse DNS lookup, but it seems simpler to just pass it with the datagram, since it is part of the user idenfifier information (and it makes it more explicit what address is used, should be consistent since its used for the account identifier. )

For datagrams that can be replayed (payment routing for example can only be replayed for the initial query, after that it depends on the routing cache), a nonce is used. The nonce is either between user (client) and server, or per account relationship in server-to-server (where both an incoming and outgoing nonce is stored, as for example the initial path finding query can be replayed and path finding is bidirectional. ) The nonce has to be higher than previous nonce, it does not need to be in order. If datagrams end up out of order and a nonce is blocked, it will not matter since the most important datagrams such as payment commitments and payments are done with the routing cache as the replay protection and not with the nonce.

We serialize and deserialize the datagram to ensure the format is well defined (and use <arpa/inet.h> to ensure correct endianess for the nonce. )

    uint8_t buffer[DATAGRAM_SIZE];

    void serialize_datagram(const Datagram* dg, uint8_t* buffer) {
        size_t offset = 0;
        buffer[offset++] = dg->command;
        memcpy(buffer + offset, dg->x_username, 32); offset += 32;
        memcpy(buffer + offset, dg->y_username, 32); offset += 32;
        memcpy(buffer + offset, dg->y_server_address, 32); offset += 32;
        memcpy(buffer + offset, dg->arguments, 256); offset += 256;
        uint32_t net_nonce = htonl(dg->nonce); // Convert to network byte order
        memcpy(buffer + offset, &net_nonce, sizeof(net_nonce)); offset += sizeof(net_nonce);
        memcpy(buffer + offset, dg->signature, 32); offset += 32;
    }
    
    void deserialize_datagram(const uint8_t* buffer, Datagram* dg) {
        size_t offset = 0;
        dg->command = buffer[offset++];
        memcpy(dg->x_username, buffer + offset, 32); offset += 32;
        memcpy(dg->y_username, buffer + offset, 32); offset += 32;
        memcpy(dg->y_server_address, buffer + offset, 32); offset += 32;
        memcpy(dg->arguments, buffer + offset, 256); offset += 256;
        uint32_t net_nonce;
        memcpy(&net_nonce, buffer + offset, sizeof(net_nonce));
        dg->nonce = ntohl(net_nonce); // Convert from network byte order to host order
        memcpy(dg->signature, buffer + offset, 32); offset += 32;
    }

And for verifying the signature, using "sha2.c":

    int verify_signature(const unsigned char *serialized_dg, const Datagram* dg) {
        size_t data_size = DATAGRAM_SIZE - SIGNATURE_SIZE;
        unsigned char data_with_key[data_size + SECRET_KEY_SIZE];

        unsigned char secret_key[SECRET_KEY_SIZE];
        load_secret_key(dg, secret_key);

        // Concatenate serialized data (excluding signature) and secret key
        memcpy(data_with_key, serialized_dg, data_size);
        memcpy(data_with_key + data_size, secret_key, SECRET_KEY_SIZE);
    
        // Compute SHA-256 hash of concatenated data
        unsigned char computed_hash[32];
        sha256(data_with_key, sizeof(data_with_key), computed_hash);
    
        // Compare computed hash with provided signature
        return memcmp(computed_hash, dg->signature, sizeof(computed_hash)) == 0;
    }

The function load_secret_key takes the datagram instance as a parameter since it needs user x and user y and connection type. The path to "secretkey.txt" differs in client and server connection types. In client, its in `datadir/client/accounts/username`, and in server it is in `datadir/server/accounts/username/peers/server_address/username/`.

### Database

A datadirectory for both client and server  (tentatively at ~/.ripple, and ~/.ripple/client for client and ~/.ripple/server for server). In server, stores a folder "accounts", that stores each account on the server in a folder with the account's name. Here there is a file "secretkey.txt" with the symmetric authorization key, and also a file "nonce.txt" with the account nonce. In each account folder, there is a folder "peers", that stores account relationships. Peers are stored under both their username and their server address, first in a folder named with the server address such as "server.xyz" (or could also be an IP address), and then in a folder under their username. In the peer folders, there is also a file "secretkey.txt", and also the files "nonce_in.txt" and "nonce_out.txt", as well as the files "trustline_in.txt" and "trustline_out.txt".

### Commands

The opcodes will be divided between client and server opcodes. 0-127 will be client opcodes and 128-255 will be server opcodes. Note that commmands also use user x and user y as parameters (every command handler receives a pointer to the Datagram instance. )

Client commands:

    0. SET_TRUSTLINE
    Value: 0x00
    Description: Sets or updates a trustline to a person.
    Arguments:
    size (64 byte)
    routable (1 bit)
    
    1. GET_TRUSTLINE
    Value: 0x01
    Description: Retrieves size of trustline to a person.
    Arguments Encoding:
    
    2. PAYMENT
    Value: 0x02
    Description: Request a payment path to user Y and make a payment if path found.
    Arguments Encoding:
    identifier (32 byte)
    amount (64 byte)

    3. RETRY_PAYMENT
    Value: 0x03
    Description: Retransmit commit or finalize commands if payment got stuck.
    Arguments Encoding:
    identifier (32 byte)

    4. GET_RECEIPT
    Value: 0x04
    Description: Get payment receipt.
    Arguments Encoding:
    identifier (32 byte)

    5. CLEAR_RECEIPTS
    Value: 0x05
    Description: Delete all receipts.
    Arguments Encoding:

Server commands:
    
    128. SYNC_TRUSTLINE
    Value: 0x80
    Description: Synchronize trustline update between two accounts
    Arguments Encoding:
    size (64 byte)

    129. FIND_PATH
    Value: 0x81
    Description: Search path.
    Arguments Encoding:
    identifier (32 byte)
    size (64 byte)
    depth (1 byte)

    130. COMMIT_PAYMENT
    Value: 0x82
    Description: Commit a payment down a path.
    Arguments Encoding:
    identifier (32 byte)

    131. FINALIZE_PAYMENT
    Value: 0x83
    Description: Finalize a payment down a path.
    Arguments Encoding:
    identifier (32 byte)

### Routing

The routing is very simple. It is practically “stateless”, no routing tables are stored, all routing is generated for each payment request. The benefit is that paths change constantly in Ripple (as trust lines fill up or credit is cleared), so a “routing table” would not reflect the true state anyway.

The path-finding optimizes for never going too deep. It is bidirectional, reducing accounts queried to 2*sqrt(unidirectional). And, it searches in increments of 1, always returning to the root before increasing the depth by 1 (the root then sends out a new request, with depth incremented by 1. ) Thus, whenever a path is found, the search ends (the root stops replying to response by incrementing request. ) Path requests use an identifier that is a simple random number, and are sent both from buyer and receiver. Whenever these “fronts” meet, a path is found, and the first path found is chosen. To enforce the “return to root before incrementing” approach, accounts should only accept queries that grow in increments of 1.

The "first path found" approximates fewest hops.

The routing is centered around caches that keep track of paths an account is involved in searching for. Accounts track when they’re currently involved in a request, and they track the depth they are at for the request. Technically, linked lists are used, and linear search. During linear search (to either find a path identifier within an account’s caches, or an account within the overall routing cache) old queries are also cleared, and accounts with no active queries are cleared.

    #define CACHE_RETENTION_SECONDS 300

    typedef struct PathCacheEntry {
        int identifier;
        uint8_t pathType; // 0 for incoming, 1 for outgoing, 2 for path found
        int depth;
        char *nextHop;
        struct PathCacheEntry *next;
    } PathCacheEntry;

    typedef struct AccountNode {
        char *accountId;
        PathCacheEntry *head;
        struct AccountNode *next;
    } AccountNode;

    AccountNode *accountCache = NULL;

### Payment

A user sends a payment request, which automatically performs the path finding and completes the payment, if possible. The server stores a receipt if payment was successful, an empty file named with the payment identifier, and the user can poll the receipt with GET_RECEIPT and the payment identifier as an argument to see if it exists as in if the payment succeeded (this allows the UDP connectionless approach. ) The user can clear old receipts with a command CLEAR_RECEIPTS. Receipts accumulate unless cleared.

Payments via more than one pathway (if "bandwidth" of one pathway is not enough) have to be done manually for now. The sender and recipient have to agree to cancel the payment (essentially make backwards payment) if they do not manage to do enough payments. If the paths have already cleared such that that there is not enough trustlines to make the full refund payment, the sender can set an "unroutable" trustline that only the receipient can use exclusively (thus it cannot be accessed for routing via the network), and then set the trustline to zero again after receiving the remaining refund. The "routable" option in SET_TRUSTLINE can be used for that (set to false. )

### Misc

It is possible to add a buffer for UDP datagrams. The buffer can also do signature and nonce validation, thus be secure against spam attacks. Although the built-in UDP buffer (at OS level) may be sufficient, and extra buffering unnecessary.

The system will be single-threaded. Instead of multiple cores, can just run multiple CPU as in multiple servers, and limit a server to what can run on one core. Multiple threads is trivial to add, one solution could be to run the equivalent of "multiple servers", each "server" in its own thread, and use an account-to-thread mapping to find which thread to route datagrams to. And then do the UDP buffer and thread-routing in the main thread. And those who prefer to run multiple computers instead, can build a "virtual endpoint" that routes to servers in a cluster of server, all under the same host address, and works analogously. Or, just run smaller scale servers with fewer accounts...
