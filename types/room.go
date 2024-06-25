package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType int

const (
	_ RoomType = iota
	Economy
	Standard
	Premium
	Deluxe
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	Size      int                `bson:"size" json:"size"`
	SeaSight  bool               `bson:"seaSight" json:"seaSight"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

type QueryRoomParams struct {
	HotelID primitive.ObjectID `json:"hotelID"`
}

func (query QueryRoomParams) ToBSON() bson.M {
	b := bson.M{}

	if query.HotelID != primitive.NilObjectID {
		b["hotelID"] = query.HotelID
	}

	return b
}
