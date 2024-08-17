# Ripple in a very simple true peer-to-peer implementation

Custom transport protocol, UDP + retransmission with counters to prevent datagrams from being replayed. Since the counters are stored permanently, the "transport protocol" behaves like a permanent connection.
