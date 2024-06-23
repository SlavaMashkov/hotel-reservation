package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type UpdateHotelParams struct {
	Name     string
	Location string
	Rooms    []primitive.ObjectID
}

func (params UpdateHotelParams) ToBSON() bson.M {
	b := bson.M{}

	if len(params.Name) > 0 {
		b["name"] = params.Name
	}
	if len(params.Location) > 0 {
		b["location"] = params.Location
	}
	if len(params.Rooms) > 0 {
		b["rooms"] = params.Rooms
	}

	return b
}
