package linkedlist

import "time"

// LinkedListNodeBase serves as a base struct for linked list nodes with a timestamp, identifier, and next pointer.
type LinkedListNodeBase struct {
    Timestamp  time.Time   // Represents when the node was last modified or created
    Identifier string      // Identifier as a string
    Next       interface{} // Pointer to the next node in the linked list
}

// Find finds a node by its identifier, removes expired nodes, and returns the found node.
func (ll *LinkedListNodeBase) Find(targetID string, timeout time.Duration) *LinkedListNodeBase {
    var prev *LinkedListNodeBase
    current := ll
    now := time.Now()

    for current != nil {
        // Check if the node has expired
        if now.Sub(current.Timestamp) > timeout {
            // Remove expired node
            if prev == nil {
                ll = current.Next.(*LinkedListNodeBase)
            } else {
                prev.Next = current.Next
            }
            // If the expired node was the target, return nil
            if current.Identifier == targetID {
                return nil
            }
        } else {
            // If the node is the target and not expired, return it
            if current.Identifier == targetID {
                return current
            }
            prev = current
        }
        current = current.Next.(*LinkedListNodeBase)
    }
    return nil
}
