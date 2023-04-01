package api

import (
	"context"
	"log"

	"github.com/ehabterra/rewards/internal/models"
	"github.com/ehabterra/rewards/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type RewardsServer struct {
	pb.UnimplementedRewardsServiceServer
	DB     *gorm.DB
	Config *Config
}

func NewGRPCRewardsService(db *gorm.DB, cfg *Config) *RewardsServer {
	return &RewardsServer{DB: db, Config: cfg}
}

func (s *RewardsServer) GetPoints(ctx context.Context, request *pb.GetPointsRequest) (*pb.GetPointsResponse, error) {
	var (
		userID = request.GetUserId()
		user   = models.User{ID: userID}
	)

	err := s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		// TODO: check if not exists send another error

		log.Printf("failed to get user data: %s", err)
		return nil, status.Errorf(codes.Internal, "method GetPoints error")
	}

	return &pb.GetPointsResponse{Points: user.Points}, nil
}

func (s *RewardsServer) AddActivity(ctx context.Context, request *pb.AddActivityRequest) (*pb.AddActivityResponse, error) {
	var (
		userID             = request.GetUserId()
		actionType         = request.GetActionType()
		user               = models.User{ID: userID}
		points     float32 = 0
		err        error
	)

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		// TODO: check if not exists send another error

		log.Printf("failed to get user data: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	switch actionType {
	case pb.ActivityActionType_ActivityInvite:
		points = s.Config.Points.Invite
	case pb.ActivityActionType_ActivityAddReview:
		points = s.Config.Points.AddReview
	default:
		log.Printf("Activity type is not correct")
		return nil, status.Errorf(codes.InvalidArgument, "Activity type is not correct")
	}

	tx := s.DB.Begin()

	err = tx.Model(&user).Update("points", user.Points+points).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to update user points: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	activity := models.Activity{
		UserID: userID,
		Type:   pb.ActivityType(actionType),
		Points: points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to create user activity: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	tx.Commit()

	return &pb.AddActivityResponse{Points: user.Points}, nil
}

func (s *RewardsServer) SendPoints(ctx context.Context, request *pb.SendPointsRequest) (*pb.SendPointsResponse, error) {
	var (
		userID      = request.GetSenderId()
		recipientID = request.GetRecipientId()
		points      = request.GetPointsAmount()
		user        = models.User{ID: userID}
		recipient   = models.User{ID: recipientID}
		err         error
	)

	if points <= 0 {
		log.Printf("Shared point should be greater than 0")
		return nil, status.Errorf(codes.OutOfRange, "shared point should be greater than 0")
	}

	if userID == recipientID {
		log.Printf("sender and recipient cannot be the same")
		return nil, status.Errorf(codes.InvalidArgument, "sender and recipient cannot be the same")
	}

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		// TODO: check if not exists send another error

		log.Printf("failed to get sender data: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	if user.Points < points {
		log.Printf("user's balance is not sufficient. balance: %f, points: %f", user.Points, points)
		return nil, status.Errorf(codes.OutOfRange, "user's balance is not sufficient. balance: %f, points: %f", user.Points, points)
	}

	err = s.DB.WithContext(ctx).Where(&recipient).Find(&recipient).Error
	if err != nil {
		log.Printf("failed to get recipient data: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	tx := s.DB.Begin()

	// Sender
	err = tx.Model(&user).Update("points", user.Points-points).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to update user points: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	activity := models.Activity{
		UserID: userID,
		Type:   pb.ActivityType_SharePoints,
		Points: -points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to create user activity: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	// Recipient
	err = tx.Model(&recipient).Update("points", recipient.Points+points).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to update recipient points: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	activity = models.Activity{
		UserID: recipientID,
		Type:   pb.ActivityType_SharePoints,
		Points: points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to create recipient activity: %s", err)
		return nil, status.Errorf(codes.Internal, "method SendPoints error")
	}

	tx.Commit()

	return &pb.SendPointsResponse{Success: true}, nil
}

func (s *RewardsServer) SpendPoints(ctx context.Context, request *pb.SpendPointsRequest) (*pb.SpendPointsResponse, error) {
	var (
		userID     = request.GetUserId()
		actionType = request.GetActionType()
		points     = request.GetPointsAmount()
		objectID   = request.GetObjectId()
		user       = models.User{ID: userID}
		err        error
	)

	err = s.DB.WithContext(ctx).Where(&user).Find(&user).Error
	if err != nil {
		// TODO: check if not exists send another error

		log.Printf("failed to get user data: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	switch actionType {
	case pb.SpendActionType_SpendPayService, pb.SpendActionType_SpendPayOtherService:
	default:
		log.Printf("Action type is not correct")
		return nil, status.Errorf(codes.InvalidArgument, "action type is not correct")
	}

	tx := s.DB.Begin()

	err = tx.Model(&user).Update("points", user.Points-points).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to update user points: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	activity := models.Activity{
		UserID:   userID,
		ObjectID: objectID,
		Type:     pb.ActivityType(actionType),
		Points:   -points,
	}

	err = tx.WithContext(ctx).Create(&activity).Error
	if err != nil {
		tx.Rollback()
		log.Printf("failed to create user activity: %s", err)
		return nil, status.Errorf(codes.Internal, "method AddActivity error")
	}

	tx.Commit()

	return &pb.SpendPointsResponse{Activity: &pb.Activity{
		Id:        activity.ID,
		UserId:    activity.UserID,
		Type:      activity.Type,
		Points:    activity.Points,
		CreatedAt: timestamppb.New(activity.CreatedAt),
	}}, nil
}
