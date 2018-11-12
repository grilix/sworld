package sworld

import (
	"errors"
	"math/rand"
	"time"
)

var (
	// ErrIncompatibleZones is when two zones are not compatible
	ErrIncompatibleZones = errors.New("Zones are not compatible")
	// ErrLowLevelStones is when both stones are level 0
	ErrLowLevelStones = errors.New("Stones are too weak")
	// ErrIncompatibleStones is when stones levels are different
	ErrIncompatibleStones = errors.New("Stones are not compatible")
)

// PortalStone is used to open a portal
type PortalStone struct {
	Level        int
	Zone         *Zone
	Duration     time.Duration
	DropInterval time.Duration
}

func (s PortalStone) minDuration(stone PortalStone) PortalStone {
	if s.Duration < stone.Duration {
		return s
	}
	return stone
}

func (s PortalStone) maxInterval(stone PortalStone) PortalStone {
	if s.DropInterval > stone.DropInterval {
		return s
	}
	return stone
}

func (s PortalStone) sortLevel(stone PortalStone) (PortalStone, PortalStone) {
	if s.Level < stone.Level {
		return s, stone
	}
	return stone, s
}

// Merge creates a new stone from two
func (s PortalStone) Merge(stone PortalStone) (PortalStone, error) {
	if s.Zone.ID != stone.Zone.ID {
		return PortalStone{}, ErrIncompatibleZones
	}

	if s.Level == stone.Level {
		return PortalStone{
			Level:        s.Level + 1,
			Zone:         s.Zone,
			Duration:     s.minDuration(stone).Duration,
			DropInterval: s.maxInterval(stone).DropInterval,
		}, nil
	}

	minLevel, maxLevel := s.sortLevel(stone)

	if minLevel.Level == 0 {
		if maxLevel.Level == 0 {
			return PortalStone{}, ErrLowLevelStones
		}

		if minLevel.Duration < maxLevel.Duration {
			// TODO: Have a better error, maybe?
			return PortalStone{}, ErrIncompatibleStones
		}
		return PortalStone{
			Level:        maxLevel.Level,
			Zone:         s.Zone,
			Duration:     maxLevel.Duration + (1 * time.Second),
			DropInterval: s.maxInterval(stone).DropInterval,
		}, nil
	}
	return PortalStone{}, ErrIncompatibleStones
}

// RandomPortalStone returns a random portal stone based on current portal
func (p Portal) RandomPortalStone() *PortalStone {
	maxDuration := int(p.PortalStone.Duration.Seconds() * 1.3)
	minDuration := int(p.PortalStone.Duration.Seconds() * 0.8)
	if maxDuration < 10 {
		maxDuration = 10
	}
	if minDuration < 10 {
		minDuration = 10
	}
	diff := maxDuration - minDuration
	var seconds int
	if diff > 2 {
		seconds = rand.Intn(maxDuration-minDuration) + minDuration
	} else {
		seconds = 10
	}
	level := rand.Intn(p.PortalStone.Level + 1)

	return &PortalStone{
		Level:    level,
		Duration: time.Duration(seconds) * time.Second,
		Zone:     p.PortalStone.Zone,
	}
}
