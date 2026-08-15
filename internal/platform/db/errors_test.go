package db

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestIsDuplicateKey(t *testing.T) {
	we := mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000, Message: "E11000 duplicate key error"}}}
	if !IsDuplicateKey(we) {
		t.Error("WriteException 11000 should be duplicate")
	}
	if IsDuplicateKey(mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11002}}}) {
		t.Error("non-11000 should not be duplicate")
	}
	ce := mongo.CommandError{Code: 11001}
	if !IsDuplicateKey(ce) {
		t.Error("CommandError 11001 should be duplicate")
	}
	if IsDuplicateKey(ce) == false {
		// sanity — keep compiler happy
	}
}
