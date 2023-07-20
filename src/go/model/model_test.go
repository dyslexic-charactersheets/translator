package model

import (
	"testing"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
)

func TestEntryID(t *testing.T) {
	entry := Entry{"Level", ""}
	log.Log("test", "ID:", entry.ID())
	// Output ID: 2698725818
}