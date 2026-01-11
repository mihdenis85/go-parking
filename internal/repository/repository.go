package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/amend-parking-backend/internal/database"
	"github.com/amend-parking-backend/internal/models"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetCountOfOccupiedSpaces(ctx context.Context) (int64, error) {
	collection := database.DB.Collection(models.ParkingSpaceLog{}.CollectionName())
	filter := bson.M{"is_active": true}
	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

func (r *Repository) GetOccupiedSpaces(ctx context.Context) ([]models.ParkingSpaceLog, error) {
	collection := database.DB.Collection(models.ParkingSpaceLog{}.CollectionName())
	filter := bson.M{"is_active": true}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var spaces []models.ParkingSpaceLog
	if err = cursor.All(ctx, &spaces); err != nil {
		return nil, err
	}

	return spaces, nil
}

func (r *Repository) GetParkingSpaceLogByPlaceNumber(ctx context.Context, placeNumber int) (*models.ParkingSpaceLog, error) {
	collection := database.DB.Collection(models.ParkingSpaceLog{}.CollectionName())
	filter := bson.M{"place_number": placeNumber, "is_active": true}

	var log models.ParkingSpaceLog
	err := collection.FindOne(ctx, filter).Decode(&log)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *Repository) AddParkingSpaceLog(ctx context.Context, log *models.ParkingSpaceLog) error {
	collection := database.DB.Collection(log.CollectionName())
	result, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}
	return nil
}

func (r *Repository) GetParkingSpaceLogsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) ([]models.ParkingSpaceLog, error) {
	collection := database.DB.Collection(models.ParkingSpaceLog{}.CollectionName())
	filter := bson.M{
		"first_name": bson.M{"$regex": "^" + firstName + "$", "$options": "i"},
		"last_name":  bson.M{"$regex": "^" + lastName + "$", "$options": "i"},
		"is_active":  true,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.ParkingSpaceLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *Repository) UpdateParkingSpaceLog(ctx context.Context, log *models.ParkingSpaceLog) error {
	collection := database.DB.Collection(log.CollectionName())
	filter := bson.M{"_id": log.ID}
	update := bson.M{"$set": log}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
