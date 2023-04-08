package types

import (
	"fmt"

	"github.com/ehabterra/rewards/internal/pb"
)

type ActionType pb.ActivityActionType

func (a ActionType) GetPoints(cfg Points) (float32, error) {
	var actionPoints = map[pb.ActivityActionType]float32{
		pb.ActivityActionType_ActivityInvite:    cfg.Invite,
		pb.ActivityActionType_ActivityAddReview: cfg.AddReview,
	}

	if points, ok := actionPoints[pb.ActivityActionType(a)]; ok {
		return points, nil
	}

	return 0, fmt.Errorf("activity type is not correct")

}
