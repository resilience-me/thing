### Tranasport control block for persistent connection

Each account, in each of its peer connections, stores a series of counters.

`rcv.nxt` the central and most important counter. Says what counter value the account will accept as the next datagram. The account will not accept a datagram that does not have its counter set to `rcv.nxt`. The sender is forced to conform to the `rcv.nxt` of the receiver, and can use polling to ensure it has the most recent value.

`snd.una` is the locally stored representation of the `rcv.nxt` at the receiving end. The sender is forced to ensure this is identical to `rcv.nxt` at the receiver, and does this by polling.

`snd.nxt` is the predicted next sequence number (one higher than the previously sent). It cannot be used until it has been ensured `rcv.nxt` at the receiver has the same value.
