package mongo

import (
	"context"
	"errors"
	"github.com/CharLottery/lottery_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LotteryRepository struct {
	collection *mongo.Collection
}

func NewLotteryRepository(db *mongo.Database) *LotteryRepository {
	return &LotteryRepository{collection: db.Collection("lotteries")}
}

func (r *LotteryRepository) CreateLottery(ctx context.Context, lottery *model.Lottery) error {
	lottery.ID = primitive.NewObjectID()
	lottery.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, lottery)
	return err
}

func (r *LotteryRepository) GetLottery(ctx context.Context, id string) (*model.Lottery, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var lottery model.Lottery
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&lottery)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &lottery, nil
}

func (r *LotteryRepository) ListLotteries(ctx context.Context) ([]*model.Lottery, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lotteries []*model.Lottery
	for cursor.Next(ctx) {
		var lottery model.Lottery
		if err := cursor.Decode(&lottery); err != nil {
			return nil, err
		}
		lotteries = append(lotteries, &lottery)
	}

	return lotteries, nil
}

func (r *LotteryRepository) DecreaseAvailableTickets(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "available_tickets": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"available_tickets": -1}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("lottery not found or no tickets left")
	}
	return nil
}
