# Ripple in a very simple true peer-to-peer implementation

Custom transport protocol, UDP + retransmission and acknowledgement with ephemeral port as "nonce". At the application layer, counters to prevent datagrams from being replayed.
