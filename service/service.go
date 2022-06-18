package service

import (
	"context"
	"github.com/pkg/errors"

	"skud"
	"skud/internal/repository"
)

type SkudService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *SkudService {
	return &SkudService{repo}
}

func (svc *SkudService) CheckAccess(ctx context.Context, readerID int64, passcardCode string) (msg string, access bool, err error) {
	employeeID, err := svc.repo.GetEmployeeIDByCode(ctx, passcardCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return skud.AccessDeniedUnknownEmployee, false, nil
		}
		return "", false, errors.Wrap(err, "failed to get employee ID")
	}
	node, err := svc.repo.GetCurrentAccessNode(ctx, employeeID)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get current access node")
	}
	nodeID, ok := node.CanReach(readerID)
	if !ok {
		return skud.AccessDeniedInaccessible, false, nil
	}
	if node.ID == nodeID {
		// person tries to exit current node
		return skud.AccessGranted, true, svc.stepUpCurrentNode(ctx, employeeID)
	}
	defer func() {
		if access {
			err = svc.repo.UpdateLastBeen(ctx, employeeID, node.ID)
		}
	}()
	node = node.GetChild(nodeID)
	node.Checks, err = svc.repo.GetAccessNodeChecks(ctx, employeeID, nodeID)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get access node required checks")
	}
	access = true
	msg = skud.AccessGranted
	switch 2*boolToInt(node.Checks.HealthCheck) - boolToInt(node.Checks.SanitaryCheck) {
	case 1: // both HealthCheck and SanitaryCheck are true
		msg = skud.AccessGrantedWithAllChecks
		access = node.Checks.HealthAccess && node.Checks.SanitaryAccess
		if !access {
			msg = getDeniedMessage(node.Checks.HealthAccess, node.Checks.SanitaryAccess)
		}
	case 2: // HealthCheck is true, SanitaryCheck is false
		msg = skud.AccessGrantedWithHealthCheck
		access = node.Checks.HealthAccess
		if !access {
			msg = skud.AccessDeniedNoHealthCheck
		}
	case -1: // HealthCheck is false, SanitaryCheck is true
		msg = skud.AccessGrantedWithSanitaryCheck
		access = node.Checks.SanitaryAccess
		if !access {
			msg = skud.AccessDeniedNoSanitaryCheck
		}
	default:
	}
	return msg, access, nil
}

func (svc *SkudService) stepUpCurrentNode(ctx context.Context, employeeID int64) error {
	transitionNode, err := svc.repo.FindLastActiveTransition(ctx, employeeID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return svc.repo.UpdateLastBeen(ctx, employeeID, transitionNode.FromNode)
	}
	return svc.repo.UpdateLastBeenToParent(ctx, employeeID)
}

func getDeniedMessage(health, sanitary bool) (msg string) {
	if health && sanitary {
		// both are true - not a "Denied" case
		return
	}
	// if one is true, another is false
	if health {
		return skud.AccessDeniedNoSanitaryCheck
	}
	if sanitary {
		return skud.AccessDeniedNoHealthCheck
	}
	// if none was true then both are false
	return skud.AccessDeniedNoAnyChecks
}

func boolToInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
