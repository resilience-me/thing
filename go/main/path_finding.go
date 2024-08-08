package main

import (
    "time"
)

type PathCacheEntry struct {
    Identifier [32]byte
    Timestamp  time.Time
    Depth      int
    NextHop    [32]byte
    Next       *PathCacheEntry
}

type AccountNode struct {
    Username     [32]byte
    Incoming     *PathCacheEntry
    Outgoing     *PathCacheEntry
    LastModified time.Time
    Next         *AccountNode
}
