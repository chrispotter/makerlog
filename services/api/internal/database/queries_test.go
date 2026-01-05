package database

import (
	"testing"
)

func TestNew(t *testing.T) {
	queries := New(nil)
	if queries == nil {
		t.Error("Expected Queries to be created")
		return
	}
	if queries.db != nil {
		t.Error("Expected db to be nil when passed nil")
	}
}
