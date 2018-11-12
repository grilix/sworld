package sworld

import (
	"log"

	"github.com/encryptio/alias"
)

type eventDropFn struct {
	fn    func(*Portal, int) *PortalEvent
	rate  float64
	level int
}

type itemDropFn struct {
	fn    func(*Portal) Item
	rate  float64
	level int
}

// Zone defines the type of enemies that will be found
// It's the base for creating a portal
type Zone struct {
	ID         string
	Name       string
	itemDrops  []itemDropFn
	eventDrops []eventDropFn
}

// NewZone initializes a zone
func NewZone(name string) *Zone {
	return &Zone{
		ID:         RandomID(16),
		Name:       name,
		itemDrops:  make([]itemDropFn, 0, 5),
		eventDrops: make([]eventDropFn, 0, 5),
	}
}

// AddItemDrop registers an item drop on a zone
func (z *Zone) AddItemDrop(minLevel int, rate float64, fn func(*Portal) Item) {
	z.itemDrops = append(z.itemDrops, itemDropFn{
		fn:    fn,
		rate:  rate,
		level: minLevel,
	})
}

// AddEventDrop registers an event drop on a zone
func (z *Zone) AddEventDrop(minLevel int, rate float64, fn func(*Portal, int) *PortalEvent) {
	z.eventDrops = append(z.eventDrops, eventDropFn{
		fn:    fn,
		rate:  rate,
		level: minLevel,
	})
}

// InitializePortal initializes portal
func (z *Zone) InitializePortal(portal *Portal) error {
	level := portal.PortalStone.Level

	items := make([]float64, 0, len(z.itemDrops))
	for _, drop := range z.itemDrops {
		if level >= drop.level {
			items = append(items, drop.rate)
		}
	}
	events := make([]float64, 0, len(z.eventDrops))
	for _, drop := range z.eventDrops {
		if level >= drop.level {
			events = append(events, drop.rate)
		}
	}

	drops, err := alias.New(items)
	if err != nil {
		return err
	}
	portal.drops = drops

	drops, err = alias.New(events)
	if err != nil {
		return err
	}
	portal.eventsRate = drops

	return nil
}

func (z *Zone) DropEvent(portal *Portal, position int) *PortalEvent {
	if portal.eventsRate == nil {
		log.Println(" ZONE: portal.drops == nil")
		return nil
	}

	event := int(portal.eventsRate.Gen(portal.seed))

	// FIXME: This is a workaround for skipping items of higher levels
	level := portal.PortalStone.Level

	for i := 0; i <= event; i++ {
		drop := z.eventDrops[i]

		if drop.level > level {
			event++
		}
	}

	return z.eventDrops[event].fn(portal, position)
}

// DropItem drops a random item from this zone
func (z *Zone) DropItem(portal *Portal) Item {
	if portal.drops == nil {
		log.Println(" ZONE: portal.drops == nil")
		return nil
	}

	item := int(portal.drops.Gen(portal.seed))

	// FIXME: This is a workaround for skipping items of higher levels
	level := portal.PortalStone.Level

	for i := 0; i <= item; i++ {
		drop := z.itemDrops[i]

		if drop.level > level {
			item++
		}
	}

	return z.itemDrops[item].fn(portal)
}
