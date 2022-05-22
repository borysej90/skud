package skud

type Checks struct {
	HealthCheck  bool
	HealthAccess bool
	// SanitaryCheck shows whether sanitary check is needed to pass this node.
	SanitaryCheck bool
	// SanitaryAccess show whether a specific person has successfully passed sanitary check.
	SanitaryAccess bool
}

type AccessNode struct {
	ID       int64
	ParentID int64
	Name     string
	Checks   Checks

	EntranceReaderID int64
	ExitReaderID     int64

	Children []*AccessNode
}

// CanReach compares readerId to physically reachable readers.
// Reachable readers are either ExitReaderID or any of Children's EntranceReaderID.
//
// i.e. CanReach returns true if readerID is equal to
// ExitReaderID or any of Children's EntranceReaderID, and respective node ID.
func (n AccessNode) CanReach(readerID int64) (int64, bool) {
	if readerID == n.ExitReaderID {
		return n.ID, true
	}
	for _, child := range n.Children {
		if child.EntranceReaderID == readerID {
			return child.EntranceReaderID, true
		}
	}
	return 0, false
}

// GetChild returns pointer to a child with ID equal to nodeID.
func (n AccessNode) GetChild(nodeID int64) *AccessNode {
	for i := range n.Children {
		if n.Children[i].ID == nodeID {
			return n.Children[i]
		}
	}
	return nil
}
