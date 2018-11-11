package sworld

import (
	"testing"
	"time"
)

var (
	zone1 = &Zone{ID: "Zone 1", Name: "Test zone 1"}
	zone2 = &Zone{ID: "Zone 2", Name: "Test zone 2"}
)

func TestMergeSameLevel(t *testing.T) {
	stone1 := PortalStone{
		Level:    1,
		Duration: 11 * time.Second,
		Zone:     zone1,
	}
	stone2 := PortalStone{
		Level:    1,
		Duration: 12 * time.Second,
		Zone:     zone1,
	}
	result, err := stone1.Merge(stone2)
	if err != nil {
		t.Fatal(err)
	}

	if result.Level != 2 {
		t.Error("Expected merged stone to be level 3, got:", result.Level)
	}
	if result.Duration != stone1.Duration {
		t.Error("Expected result to have 11s duration, got:", result.Duration)
	}
	if result.Zone.ID != stone1.Zone.ID {
		t.Error("Expected stone to be for same zone, got:", result.Zone)
	}

	result, err = stone2.Merge(stone1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Level != 2 {
		t.Error("Expected merged stone to be level 3, got:", result.Level)
	}
	if result.Duration != stone1.Duration {
		t.Error("Expected result to have 11s duration, got:", result.Duration)
	}
	if result.Zone.ID != stone1.Zone.ID {
		t.Error("Expected stone to be for same zone, got:", result.Zone)
	}
}

func TestMergeLevelZero(t *testing.T) {
	stone1 := PortalStone{
		Level:    0,
		Duration: 10 * time.Second,
		Zone:     zone1,
	}
	stone2 := PortalStone{
		Level:    3,
		Duration: 15 * time.Second,
		Zone:     zone1,
	}
	stone3 := PortalStone{
		Level:    0,
		Duration: 20 * time.Second,
		Zone:     zone1,
	}
	result, err := stone2.Merge(stone1)
	if err == nil {
		// Merging a stone with a stone of level 0 requires the level 0 stone
		// to have a longer duration
		t.Error("Expected merge to fail, got ", result)
	}
	result, err = stone2.Merge(stone3)
	if err != nil {
		t.Fatal("Expected merge to be successful, got:", err)
	}

	if result.Level != 3 {
		t.Error("Expected result to be level 3, got:", result.Level)
	}

	if result.Duration != (stone2.Duration + (1 * time.Second)) {
		t.Error("Expected duration to have increased by 1 second, got:", result)
	}
}

func TestMergeDifferentLevels(t *testing.T) {
	stone1 := PortalStone{
		Level:    2,
		Duration: 10 * time.Second,
		Zone:     zone1,
	}
	stone2 := PortalStone{
		Level:    3,
		Duration: 15 * time.Second,
		Zone:     zone1,
	}
	result, err := stone1.Merge(stone2)
	if err == nil {
		t.Error("Expected merge to fail, got ", result)
	}
}
