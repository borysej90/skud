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
}

func (r accessNodeRecord) toAccessNode() *skud.AccessNode {
	return &skud.AccessNode{
		ID:               r.ID.Int64,
		Name:             r.Name,
		ParentID:         r.ParentID.Int64,
		EntranceReaderID: r.EntranceReaderID,
		ExitReaderID:     r.ExitReaderID.Int64,
	}
}
