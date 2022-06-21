package sql

import (
	"database/sql"

	"skud"
)

type accessNodeRecord struct {
	ID               sql.NullInt64 `db:"id"`
	ParentID         sql.NullInt64 `db:"parent_id"`
	Name             string        `db:"name"`
	EntranceReaderID int64         `db:"entrance_reader"`
	ExitReaderID     sql.NullInt64 `db:"exit_reader"`
	TransitiveTo     sql.NullInt64 `db:"transitive_to"`
}

func (r accessNodeRecord) toAccessNode() *skud.AccessNode {
	return &skud.AccessNode{
		ID:               r.ID.Int64,
		Name:             r.Name,
		ParentID:         r.ParentID.Int64,
		TransitiveTo:     r.TransitiveTo.Int64,
		EntranceReaderID: r.EntranceReaderID,
		ExitReaderID:     r.ExitReaderID.Int64,
	}
}

type transitionNodeRecord struct {
	ID       int64 `db:"id"`
	From     int64 `db:"from_node"`
	To       int64 `db:"to_node"`
	ParentID int64 `db:"parent_id"`
}

func (r transitionNodeRecord) toTransitionNode() *skud.TransitionNode {
	return &skud.TransitionNode{
		ID:         r.ID,
		FromNode:   r.From,
		ToNode:     r.To,
		ParentNode: r.ParentID,
	}
}
