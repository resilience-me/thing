### Tranasport control block for persistent connection

Each account, in each of its peer connections, stores a series of counters.

`rcv.nxt` the central and most important counter. Says what counter value the account will accept as the next datagram. The account will not accept a datagram that does not have its counter set to rcv.nxt.
