package main

import (
    "time"
)

type PathEntry struct {
    Identifier [32]byte
    Timestamp  time.Time
    Depth      int
    NextHop    [32]byte
    Next       *PathEntry
}

type AccountNode struct {
    Username     [32]byte
    Incoming     *PathEntry
    Outgoing     *PathEntry
    LastModified time.Time
    Next         *AccountNode
}
