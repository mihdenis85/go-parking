package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParkingSpaceLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"507f1f77bcf86cd799439011"`
	LogID        string             `bson:"log_id" json:"log_id" example:"log-123"`
	PlaceNumber  int                `bson:"place_number" json:"place_number" example:"1"`
	FirstName    string             `bson:"first_name" json:"first_name" example:"Иван"`
	LastName     string             `bson:"last_name" json:"last_name" example:"Иванов"`
	CarMake      string             `bson:"car_make" json:"car_make" example:"Toyota"`
	LicensePlate string             `bson:"license_plate" json:"license_plate" example:"А123БВ777"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at" example:"2024-01-01T12:00:00Z"`
	IsActive     bool               `bson:"is_active" json:"is_active" example:"true"`
	FreeUpTime   *time.Time         `bson:"free_up_time,omitempty" json:"free_up_time,omitempty" example:"2024-01-01T14:00:00Z"`
}

func (p ParkingSpaceLog) CollectionName() string {
	return "parking_space_logs"
}
