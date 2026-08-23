package mongorepo

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func extractDuplicateField(err error) (field string, value interface{}, ok bool) {
	var we *mongo.WriteException
	if !errors.As(err, &we) {
		return "", nil, false
	}

	for _, e := range we.WriteErrors {
		if e.Code != 11000 {
			continue
		}

		var details struct {
			KeyValue bson.Raw `bson:"keyValue"`
		}
		if err := json.Unmarshal(e.Raw, &details); err != nil {
			continue
		}

		elems, err := details.KeyValue.Elements()
		if err != nil || len(elems) == 0 {
			continue
		}

		return elems[0].Key(), elems[0].Value(), true
	}
	return "", nil, false
}
