package linkedlist

import (
    "time"
    "ripple/config"
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

func (bl *BaseList) Add(newNode *BaseNode) {
    newNode.Timestamp = time.Now() // Set the current time as the timestamp
    newNode.Next = bl.head         // Point the new node's next to the current head
    bl.head = newNode              // Update the head to be the new node
}

// Find finds a node by its identifier, removes expired nodes, and returns the found node.
func (bl *BaseList) Find(identifier string) *BaseNode {
    var prev *BaseNode
    current := bl.head
    now := time.Now()

    for current != nil {
        // Check if the node has expired
        if now.Sub(current.Timestamp) > config.PathFindingTimeout {
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

// Remove deletes a node by its identifier from the list and returns true if the node was found and removed.
func (bl *BaseList) Remove(identifier string) bool {
    // Start by traversing the list to find the node and its previous node
    var prev *BaseNode
    current := bl.head

    for current != nil {
        if current.Identifier == identifier {
            // If the node to remove is the head
            if prev == nil {
                bl.head = current.Next
            } else {
                // Bypass the node to remove it
                prev.Next = current.Next
            }

            return true
        }
        prev = current
        current = current.Next
    }

    // Return false if the node was not found
    return false
}
