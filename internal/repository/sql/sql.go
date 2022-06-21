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
	err := r.db.GetContext(ctx, &record, r.db.Rebind(stmt), employeeID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, translateDBErr(err)
	}
	node := record.toAccessNode()
	stmt = `
SELECT an.id AS id, an.name AS name, parent_id, entrance_reader, exit_reader, to_node AS transitive_to FROM access_nodes AS an
LEFT JOIN transition_nodes AS tn ON tn.from_node = an.id
WHERE parent_id`
	args := make([]interface{}, 0, 1)
	if record.ID.Valid {
		stmt += " = ?"
		args = append(args, record.ID)
	} else {
		stmt += " IS NULL"
	}
	records := make([]*accessNodeRecord, 0)
	if err = r.db.SelectContext(ctx, &records, r.db.Rebind(stmt), args...); err != nil {
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
SELECT tn.id AS id, from_node, to_node, an.parent_id AS parent_id FROM transition_nodes AS tn
JOIN transits AS t ON t.transition_node_id = tn.id
JOIN employees AS e ON e.id = t.employee_id
JOIN access_nodes AS an ON an.id = tn.from_node
WHERE e.id = ? AND e.last_been = tn.to_node`
	var record transitionNodeRecord
	if err := r.db.GetContext(ctx, &record, r.db.Rebind(stmt), employeeID); err != nil {
		return nil, translateDBErr(err)
	}
	return record.toTransitionNode(), nil
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

func (r *Repo) TransitForward(ctx context.Context, employeeID, fromNode int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // it's no-op if committed
	stmt := `
UPDATE employees SET last_been = (
	SELECT to_node FROM transition_nodes WHERE from_node = ?
) WHERE id = ?`
	if _, err = tx.ExecContext(ctx, tx.Rebind(stmt), fromNode, employeeID); err != nil {
		return translateDBErr(err)
	}
	stmt = "INSERT INTO transits VALUES (?, (SELECT id FROM transition_nodes WHERE from_node = ?))"
	if _, err = tx.ExecContext(ctx, tx.Rebind(stmt), employeeID, fromNode); err != nil {
		return translateDBErr(err)
	}
	return tx.Commit()
}

func (r *Repo) TransitBackward(ctx context.Context, employeeID, transitionNodeID int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // it's no-op if committed
	stmt := `
UPDATE employees AS e
JOIN transits AS t ON t.employee_id = e.id
JOIN transition_nodes AS tn ON tn.id = t.transition_node_id
JOIN access_nodes AS an ON an.id = tn.from_node
SET last_been = parent_id
WHERE e.id = ? AND last_been = to_node`
	if _, err = tx.ExecContext(ctx, tx.Rebind(stmt), employeeID); err != nil {
		return translateDBErr(err)
	}
	stmt = "DELETE FROM transits WHERE employee_id = ? AND transition_node_id = ?"
	if _, err = tx.ExecContext(ctx, tx.Rebind(stmt), employeeID, transitionNodeID); err != nil {
		return translateDBErr(err)
	}
	return tx.Commit()
}

func translateDBErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}
