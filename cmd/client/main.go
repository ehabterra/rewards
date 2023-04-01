package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ehabterra/rewards/internal/pb"
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
	userID      = flag.Int("user_id", defaultUserID, "User ID")
	recipientID = flag.Int("recipient_id", defaultRecipientID, "Recipient ID")
	points      = flag.Float64("points", defaultPoints, "Shared points")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRewardsServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Add activity
	r, err := c.AddActivity(ctx, &pb.AddActivityRequest{
		UserId:     int32(*userID),
		ActionType: pb.ActivityActionType_ActivityAddReview,
	})
	if err != nil {
		log.Fatalf("failed to update: %v", err)
	}
	log.Printf("balance: %f", r.Points)

	// Send to friend
	_, err = c.SendPoints(ctx, &pb.SendPointsRequest{
		SenderId:     int32(*userID),
		RecipientId:  int32(*recipientID),
		PointsAmount: float32(*points),
	})
	if err != nil {
		log.Fatalf("failed to share points: %v", err)
	}
	log.Printf("successfully shared points: %f", *points)

	// Get user points
	userPoints, err := c.GetPoints(ctx, &pb.GetPointsRequest{
		UserId: int32(*userID),
	})
	if err != nil {
		log.Fatalf("failed to get points: %v", err)
	}
	log.Printf("user points: %f", userPoints.Points)

	// Get user points
	recipientPoints, err := c.GetPoints(ctx, &pb.GetPointsRequest{
		UserId: int32(*recipientID),
	})
	if err != nil {
		log.Fatalf("failed to get points: %v", err)
	}
	log.Printf("recipient points: %f", recipientPoints.Points)
}
