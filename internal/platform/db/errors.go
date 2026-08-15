package db

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// IsDuplicateKey reports whether err is a Mongo unique-index violation (E11000).
func IsDuplicateKey(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 || e.Code == 11001 {
				return true
			}
		}
	}
	var ce mongo.CommandError
	if errors.As(err, &ce) {
		if ce.Code == 11000 || ce.Code == 11001 {
			return true
		}
	}
	return false
}
