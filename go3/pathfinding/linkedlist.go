package pathfinding

import (
    "time"
)

// BaseNode serves as a base struct for linked list nodes with a timestamp, identifier, and next pointer.
type BaseNode struct {
    Timestamp  time.Time   // Represents when the node was last modified or created
    Identifier string      // Identifier as a string
    Next       *BaseNode   // Pointer to the next node in the linked list
}

// BaseList manages a linked list of BaseNode elements.
type BaseList struct {
    head *BaseNode // Pointer to the head of the linked list
}

// Add creates a new node with the provided identifier and inserts it at the head of the list.
func (bl *BaseList) Add(identifier string) *BaseNode {
    newNode := &BaseNode{
        Timestamp:  time.Now(),
        Identifier: identifier,
        Next:       bl.head, // Point the new node's next to the current head
    }
    bl.head = newNode // Update the head to be the new node
    return newNode    // Return the newly created node
}

// Find finds a node by its identifier, removes expired nodes, and returns the found node.
func (bl *BaseList) Find(identifier string, timeout time.Duration) *BaseNode {
    var prev *BaseNode
    current := bl.head
    now := time.Now()

    for current != nil {
        // Check if the node has expired
        if now.Sub(current.Timestamp) > timeout {
            // Remove expired node
            if prev == nil {
                bl.head = current.Next
            } else {
                prev.Next = current.Next
            }
            // If the expired node was the target, return nil
            if current.Identifier == identifier {
                return nil
            }
        } else {
            // If the node is the target and not expired, return it
            if current.Identifier == identifier {
                return current
            }
            prev = current
        }
        current = current.Next
    }
    return nil
}
