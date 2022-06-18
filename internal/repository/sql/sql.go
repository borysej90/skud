package sql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"skud"
	"skud/internal/repository"
)

var _ repository.Repository = (*Repo)(nil)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) GetEmployeeIDByCode(ctx context.Context, code string) (int64, error) {
	stmt := "SELECT id FROM employees WHERE card = ? AND active = 1"
	var ret int64
	return ret, translateDBErr(r.db.GetContext(ctx, &ret, r.db.Rebind(stmt), code))
}

func (r *Repo) GetCurrentAccessNode(ctx context.Context, employeeID int64) (*skud.AccessNode, error) {
	stmt := `
	SELECT * FROM access_nodes
	WHERE id = (SELECT last_been FROM employees WHERE id = ?)`
	var record accessNodeRecord
	if err := r.db.GetContext(ctx, &record, r.db.Rebind(stmt), employeeID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, translateDBErr(err)
	}
	node := record.toAccessNode()
	stmt = `SELECT * FROM access_nodes WHERE parent_id IS NULL`
	args := make([]interface{}, 0, 1)
	if record.ID.Valid {
		args = append(args, record.ID)
		stmt = `SELECT * FROM access_nodes WHERE parent_id = ?`
	}
	records := make([]*accessNodeRecord, 0)
	if err := r.db.SelectContext(ctx, &records, r.db.Rebind(stmt), args...); err != nil {
		return nil, translateDBErr(err)
	}
	node.Children = make([]*skud.AccessNode, len(records))
	for i, rec := range records {
		node.Children[i] = rec.toAccessNode()
	}
	return node, nil
}

func (r *Repo) GetAccessNodeChecks(ctx context.Context, employeeID, nodeID int64) (skud.Checks, error) {
	stmt := `
SELECT
       COALESCE(health_check, 0),
       CASE WHEN health_check THEN (
		   SELECT COUNT(id) > 0 FROM health_checks WHERE employee_id = ? AND until > CURRENT_TIMESTAMP
	   ) ELSE 1 END AS health_access,
       COALESCE(sanitary_check, 0),
       1 AS sanitary_access -- TODO: handle sanitary_check (there is no tables with data yet)
FROM (
	SELECT AVG(health_check) = 1 AS health_check, AVG(sanitary_check) = 1 AS sanitary_check
	FROM permissions
	JOIN members USING (group_id)
	WHERE node_id = ? AND employee_id = ?
) AS sub`
	var healthCheck, healthAccess, sanitaryCheck, sanitaryAccess int
	row := r.db.QueryRowxContext(ctx, r.db.Rebind(stmt), employeeID, nodeID, employeeID)
	if err := row.Scan(&healthCheck, &healthAccess, &sanitaryCheck, &sanitaryAccess); err != nil {
		return skud.Checks{}, translateDBErr(err)
	}
	return skud.Checks{
		HealthCheck:    healthCheck == 1,
		HealthAccess:   healthAccess == 1,
		SanitaryCheck:  sanitaryCheck == 1,
		SanitaryAccess: sanitaryAccess == 1,
	}, nil
}

func (r *Repo) FindLastActiveTransition(ctx context.Context, employeeID int64) (*skud.TransitionNode, error) {
	stmt := `
SELECT tn.id, from_node, to_node FROM transition_nodes AS tn
JOIN transits AS t on t.transition_node_id = tn.id
JOIN employees AS e on e.id = t.employee_id
WHERE e.id = ? AND e.last_been = tn.to_node`
	var record transitionNodeRecord
	return record.toTransitionNode(), translateDBErr(r.db.GetContext(ctx, &record, r.db.Rebind(stmt), employeeID))
}

func (r *Repo) UpdateLastBeen(ctx context.Context, employeeID, nodeID int64) error {
	stmt := "UPDATE employees SET last_been = ? WHERE id = ?"
	lastBeen := sql.NullInt64{Int64: nodeID, Valid: nodeID != 0}
	_, err := r.db.ExecContext(ctx, r.db.Rebind(stmt), lastBeen, employeeID)
	return translateDBErr(err)
}

func (r *Repo) UpdateLastBeenToParent(ctx context.Context, employeeID int64) error {
	stmt := `
UPDATE employees AS e
JOIN access_nodes an ON an.id = e.last_been
SET e.last_been = an.parent_id
WHERE e.id = ?`
	_, err := r.db.ExecContext(ctx, r.db.Rebind(stmt), employeeID)
	return translateDBErr(err)
}

func translateDBErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}
