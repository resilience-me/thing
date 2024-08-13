package pathfinding

// FindOrAdd looks for an existing account and returns it, or adds a new one if it does not exist.
func (pm *PathManager) FindOrAdd(username string) *AccountNode {
    // Attempt to find the existing node
    existingNode := pm.Find(username)
    if existingNode != nil {
        return existingNode
    }

    // Only reach here if no existing node was found; add a new one
    return pm.Add(username)
}
