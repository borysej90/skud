package repository

import (
	"context"
	"github.com/pkg/errors"

	"skud"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	// GetEmployeeIDByCode returns person ID by pass code.
	GetEmployeeIDByCode(ctx context.Context, code string) (int64, error)

	// GetCurrentAccessNode returns the latest access node that employee has accessed with all direct
	// children that can be physically accessed.
	GetCurrentAccessNode(ctx context.Context, employeeID int64) (*skud.AccessNode, error)

	// GetAccessNodeChecks returns required checks that have to be passed by employee
	// before entering node with ID equal to nodeID and actual checks results.
	GetAccessNodeChecks(ctx context.Context, employeeID, nodeID int64) (skud.Checks, error)

	// FindLastActiveTransition returns the last employee's transition to their current node.
	//
	// If none is found, ErrNotFound is returned.
	FindLastActiveTransition(ctx context.Context, employeeID int64) (*skud.TransitionNode, error)

	UpdateLastBeen(ctx context.Context, employeeID, nodeID int64) error

	UpdateLastBeenToParent(ctx context.Context, employeeID int64) error
}
