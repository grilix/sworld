package sworld

import (
	"time"
)

// DropRate will hold information about the rates for drops
type DropRate struct {
	Gold     float64
	Enemy    float64
	Item     float64
	Stone    float64
	Material float64
	Nothing  float64
}

// Zone defines the type of enemies that will be found.
// It's the base for creating a portal.
type Zone struct {
	ID       string
	Name     string
	DropRate DropRate
}

// PortalStone is used to open a portal
type PortalStone struct {
	Level        int
	Zone         *Zone
	Duration     time.Duration
	DropInterval time.Duration
}

// Portal is an instance of a Zone, where players can teleport to.
type Portal struct {
	ID          string
	PortalStone PortalStone
	User        *User
	Character   *Character
	IsOpen      bool
}
