package sworld

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomItemEvent(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)

	portal := Portal{}
	event := portal.randomItemEvent(seed)
	if event.Item == nil {
		t.Fatal("Expected event to be Item")
	}
	_, ok := event.Item.(*PortalStone)
	if !ok {
		t.Fatal("Expected item to be a stone")
	}
}
