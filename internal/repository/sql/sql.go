package sql

import (
	"context"
	"github.com/jmoiron/sqlx"

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
	if err := r.db.GetContext(ctx, &ret, r.db.Rebind(stmt), code); err != nil {
		return 0, err
	}
	return ret, nil
}

func (r *Repo) GetCurrentAccessNode(ctx context.Context, employeeID int64) (*skud.AccessNode, error) {
	stmt := `
	SELECT * FROM access_nodes
	WHERE id = (SELECT node_id FROM employee_current_nodes WHERE employee_id = ? LIMIT 1)`
	var record accessNodeRecord
	if err := r.db.GetContext(ctx, &record, r.db.Rebind(stmt), employeeID); err != nil {
		return nil, err
	}
	node := record.toAccessNode()
	stmt = `SELECT * FROM access_nodes WHERE parent_id = ?`
	records := make([]*accessNodeRecord, 0)
	if err := r.db.SelectContext(ctx, &records, r.db.Rebind(stmt), node.ID); err != nil {
		return nil, err
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
       health_check,
       CASE WHEN health_check THEN (
		   SELECT COUNT(id) > 0 FROM health_checks WHERE employee_id = ? AND until > CURRENT_TIMESTAMP
	   ) ELSE 1 END AS health_access,
       sanitary_check,
       1 AS sanitary_access -- TODO: handle sanitary_check (there is no tables with data yet)
FROM (
	SELECT AVG(health_check) = 1 AS health_check, AVG(sanitary_check) = 1 AS sanitary_check
	FROM permissions
	JOIN members USING (group_id)
	WHERE node_id = ? AND employee_id = ?
) AS sub`
	dest := make(map[string]interface{})
	if err := r.db.GetContext(ctx, &dest, r.db.Rebind(stmt), nodeID, employeeID); err != nil {
		return skud.Checks{}, err
	}
	return skud.Checks{
		HealthCheck:    dest["health_check"].(int) == 1,
		HealthAccess:   dest["health_access"].(int) == 1,
		SanitaryCheck:  dest["sanitary_check"].(int) == 1,
		SanitaryAccess: dest["sanitary_access"].(int) == 1,
	}, nil
}
