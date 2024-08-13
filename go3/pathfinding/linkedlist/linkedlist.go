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

// FindParent searches for the parent of a node by its identifier, removes expired nodes while traversing,
// and returns the parent node. If the node is the head, parent will be nil.
func (bl *BaseList) FindParent(identifier string) *BaseNode {
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
            // If the expired node was the target, return nil (since it's been removed)
            if current.Identifier == identifier {
                return nil
            }
        } else {
            // If the node is the target and not expired, return its parent
            if current.Identifier == identifier {
                return prev
            }
            prev = current
        }
        current = current.Next
    }
    return nil
}

// Find searches for a node by its identifier, removes expired nodes while traversing,
// and returns the found node.
func (bl *BaseList) Find(identifier string) *BaseNode {
    // First, check if the list is empty
    if bl.head == nil {
        return nil
    }

    // If the head is the target node, return it directly
    if bl.head.Identifier == identifier {
        return bl.head
    }

    // Otherwise, find the parent node of the target
    parent := bl.FindParent(identifier)

    // If the parent is found, return its next node (which should be the target node)
    if parent != nil {
        return parent.Next
    }

    // If the node wasn't found, return nil
    return nil
}

// Remove searches for a node by its identifier and removes it from the list.
func (bl *BaseList) Remove(identifier string) bool {
    // Check if the list is empty
    if bl.head == nil {
        return false
    }

    // Special case: if the head is the node to be removed
    if bl.head.Identifier == identifier {
        bl.head = bl.head.Next // Move the head to the next node
        return true
    }

    // Find the parent of the node to be removed
    parent := bl.FindParent(identifier)

    // If the parent is found, remove the target node
    if parent != nil {
        parent.Next = parent.Next.Next // Bypass the node to remove it
        return true
    }

    // Node not found
    return false
}

