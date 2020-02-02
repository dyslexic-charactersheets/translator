package model

import (
	"testing"
	"../log"
)

func TestEntryID(t *testing.T) {
	entry := Entry{"Level", ""}
	log.Log("test", "ID:", entry.ID())
	// Output ID: 2698725818
}