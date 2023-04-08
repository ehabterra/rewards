package api

import (
	"context"
	"errors"

	"github.com/ehabterra/rewards/internal/models"
	"github.com/ehabterra/rewards/internal/pb"
	"github.com/ehabterra/rewards/internal/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type RewardsServer struct {
	pb.UnimplementedRewardsServiceServer
	DB     *gorm.DB
	Config *types.Config
}

func NewGRPCRewardsService(db *gorm.DB, cfg *types.Config) *RewardsServer {
	return &RewardsServer{DB: db, Config: cfg}
}

func (s *RewardsServer) GetPoints(ctx context.Context, request *pb.GetPointsRequest) (*pb.GetPointsResponse, error) {
	var (
		userID = request.GetUserId()
		user   = models.User{ID: userID}
	)

	err := s.DB.WithContext(ctx).Where(&user).Select("points").Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithContext(ctx).WithError(err).Error("user is not found")
			return nil, status.Errorf(codes.Internal, "user is not found")
		}

		log.WithContext(ctx).WithError(err).Error("failed to get user data")
		return nil, status.Errorf(codes.Internal, "failed to get user data")
	}

	return &pb.GetPointsResponse{Points: user.Points}, nil
}

func (s *RewardsServer) AddActivity(ctx context.Context, request *pb.AddActivityRequest) (*pb.AddActivityResponse, error) {
	var (
		userID             = request.GetUserId()
		actionType         = types.ActionType(request.GetActionType())
		user               = models.User{ID: userID}
		points     float32 = 0
		panicked           = true
		err        error
	)

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithContext(ctx).WithError(err).Error("user is not found")
			return nil, status.Errorf(codes.Internal, "user is not found")
		}

		log.WithContext(ctx).WithError(err).Error("failed to get user data")
		return nil, status.Errorf(codes.Internal, "failed to get user data")
	}

	// Get points based on config and action type
	points, err = actionType.GetPoints(s.Config.Points)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("activity type is not correct")
		return nil, status.Errorf(codes.InvalidArgument, "activity type is not correct")
	}

	tx := s.DB.Begin()

	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	err = tx.WithContext(ctx).Model(&user).Update("points", gorm.Expr("points + ?", points)).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to update user points")
		return nil, status.Errorf(codes.Internal, "failed to update user points")
	}

	activity := models.Activity{
		UserID: userID,
		Type:   pb.ActivityType(actionType),
		Points: points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to create user activity")
		return nil, status.Errorf(codes.Internal, "failed to create user activity")
	}

	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		panicked = false
		log.WithContext(ctx).WithError(err).Error("failed to commit user activity")
		return nil, status.Errorf(codes.Internal, "failed to commit user activity")
	}

	userPoints, err := s.GetPoints(ctx, &pb.GetPointsRequest{UserId: user.ID})
	if err != nil {
		panicked = false
		log.WithContext(ctx).WithError(err).Warning("failed to get user points")
		return nil, status.Errorf(codes.Internal, "failed to get user points")
	}

	panicked = false

	return &pb.AddActivityResponse{Points: userPoints.Points}, nil
}

func (s *RewardsServer) SendPoints(ctx context.Context, request *pb.SendPointsRequest) (*pb.SendPointsResponse, error) {
	var (
		userID      = request.GetSenderId()
		recipientID = request.GetRecipientId()
		points      = request.GetPointsAmount()
		user        = models.User{ID: userID}
		recipient   = models.User{ID: recipientID}
		panicked    = true
		err         error
	)

	if points <= 0 {
		log.WithContext(ctx).Error("Shared point should be greater than 0")
		return nil, status.Errorf(codes.OutOfRange, "shared point should be greater than 0")
	}

	if userID == recipientID {
		log.WithContext(ctx).Error("sender and recipient cannot be the same")
		return nil, status.Errorf(codes.InvalidArgument, "sender and recipient cannot be the same")
	}

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithContext(ctx).WithError(err).Error("user is not found")
			return nil, status.Errorf(codes.Internal, "user is not found")
		}

		log.WithContext(ctx).WithError(err).Error("failed to get sender data")
		return nil, status.Errorf(codes.Internal, "failed to get sender data")
	}

	if user.Points < points {
		log.WithContext(ctx).Errorf("user's balance is not sufficient. balance: %f, points: %f", user.Points, points)
		return nil, status.Errorf(codes.OutOfRange, "user's balance is not sufficient. balance: %f, points: %f", user.Points, points)
	}

	err = s.DB.WithContext(ctx).Where(&recipient).Find(&recipient).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get recipient data")
		return nil, status.Errorf(codes.Internal, "failed to get recipient data")
	}

	tx := s.DB.Begin()

	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	// Sender
	err = tx.WithContext(ctx).Model(&user).Update("points", gorm.Expr("points - ?", points)).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to update user points")
		return nil, status.Errorf(codes.Internal, "failed to update user points")
	}

	activity := models.Activity{
		UserID: userID,
		Type:   pb.ActivityType_SharePoints,
		Points: -points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to create user activity")
		return nil, status.Errorf(codes.Internal, "failed to create user activity")
	}

	// Recipient
	err = tx.WithContext(ctx).Model(&recipient).Update("points", gorm.Expr("points + ?", points)).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to update recipient points")
		return nil, status.Errorf(codes.Internal, "failed to update recipient points")
	}

	activity = models.Activity{
		UserID: recipientID,
		Type:   pb.ActivityType_SharePoints,
		Points: points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to create recipient activity")
		return nil, status.Errorf(codes.Internal, "failed to create recipient activity")
	}

	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		panicked = false
		log.WithContext(ctx).WithError(err).Error("failed to commit user points")
		return nil, status.Errorf(codes.Internal, "failed to commit user points")
	}

	panicked = false

	return &pb.SendPointsResponse{Success: true}, nil
}

func (s *RewardsServer) SpendPoints(ctx context.Context, request *pb.SpendPointsRequest) (*pb.SpendPointsResponse, error) {
	var (
		userID     = request.GetUserId()
		actionType = request.GetActionType()
		points     = request.GetPointsAmount()
		objectID   = request.GetObjectId()
		user       = models.User{ID: userID}
		panicked   = true
		err        error
	)

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithContext(ctx).WithError(err).Error("user is not found")
			return nil, status.Errorf(codes.Internal, "user is not found")
		}

		log.WithContext(ctx).WithError(err).Error("failed to get user data")
		return nil, status.Errorf(codes.Internal, "failed to get user data")
	}

	// Check if selected action is correct
	switch actionType {
	case pb.SpendActionType_SpendPayService, pb.SpendActionType_SpendPayOtherService:
	default:
		log.WithContext(ctx).Error("Action type is not correct")
		return nil, status.Errorf(codes.InvalidArgument, "action type is not correct")
	}

	tx := s.DB.Begin()

	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	err = tx.WithContext(ctx).Model(&user).Update("points", gorm.Expr("points - ?", points)).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to update user points")
		return nil, status.Errorf(codes.Internal, "failed to update user points")
	}

	activity := models.Activity{
		UserID:   userID,
		ObjectID: objectID,
		Type:     pb.ActivityType(actionType),
		Points:   -points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to create user activity")
		return nil, status.Errorf(codes.Internal, "failed to create user activity")
	}

	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		panicked = false
		log.WithContext(ctx).WithError(err).Error("failed to commit user activity")
		return nil, status.Errorf(codes.Internal, "failed to commit user activity")
	}

	panicked = false

	return &pb.SpendPointsResponse{Activity: &pb.Activity{
		Id:        activity.ID,
		UserId:    activity.UserID,
		Type:      activity.Type,
		Points:    activity.Points,
		CreatedAt: timestamppb.New(activity.CreatedAt),
	}}, nil
}
