package sql

import (
	"database/sql"

	"skud"
)

type accessNodeRecord struct {
	ID               int64          `db:"id"`
	ParentID         *sql.NullInt64 `db:"parent_id"`
	Name             string         `db:"name"`
	EntranceReaderID int64          `db:"entrance_reader_id"`
	ExitReaderID     *sql.NullInt64 `db:"exit_reader_id"`
}

func (r accessNodeRecord) toAccessNode() *skud.AccessNode {
	return &skud.AccessNode{
		ID:               r.ID,
		Name:             r.Name,
		ParentID:         r.ParentID.Int64,
		EntranceReaderID: r.EntranceReaderID,
		ExitReaderID:     r.ExitReaderID.Int64,
	}
}
