package main

import (
	"context"
	"flag"
	"time"

	"github.com/ehabterra/rewards/internal/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultUserID      = 1
	defaultRecipientID = 2
	defaultPoints      = 5.5
)

var (
	addr        = flag.String("addr", "localhost:8000", "the address to connect to")
	userID      = flag.Uint("user_id", defaultUserID, "User ID")
	recipientID = flag.Uint("recipient_id", defaultRecipientID, "Recipient ID")
	points      = flag.Float64("points", defaultPoints, "Shared points")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).Fatal("did not connect")
	}
	defer conn.Close()
	c := pb.NewRewardsServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Add activity
	r, err := c.AddActivity(ctx, &pb.AddActivityRequest{
		UserId:     uint32(*userID),
		ActionType: pb.ActivityActionType_ActivityAddReview,
	})
	if err != nil {
		log.WithError(err).Fatal("failed to update")
	}
	log.Debugf("balance: %f", r.Points)

	// Send to friend
	_, err = c.SendPoints(ctx, &pb.SendPointsRequest{
		SenderId:     uint32(*userID),
		RecipientId:  uint32(*recipientID),
		PointsAmount: float32(*points),
	})
	if err != nil {
		log.WithError(err).Fatal("failed to share points")
	}
	log.Debugf("successfully shared points: %f", *points)

	// Get user points
	userPoints, err := c.GetPoints(ctx, &pb.GetPointsRequest{
		UserId: uint32(*userID),
	})
	if err != nil {
		log.WithError(err).Fatal("failed to get points")
	}
	log.Debugf("user points: %f", userPoints.Points)

	// Get user points
	recipientPoints, err := c.GetPoints(ctx, &pb.GetPointsRequest{
		UserId: uint32(*recipientID),
	})
	if err != nil {
		log.WithError(err).Fatal("failed to get points")
	}
	log.Debugf("recipient points: %f", recipientPoints.Points)
}
