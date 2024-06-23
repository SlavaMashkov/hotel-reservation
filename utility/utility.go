package utility

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func IDToMongoOID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("could not convert to ObjectID", slog.String("id", id))
		return primitive.NilObjectID, err
	}

	return oid, nil
}
