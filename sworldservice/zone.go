package sworldservice

import (
	"time"

	"github.com/grilix/sworld/sworld"
)

// TODO: This is me testing it
func randomPowerStone(portal *sworld.Portal) sworld.Item {
	return &sworld.PortalStone{
		Level:    10,
		Duration: 10 * time.Minute,
		Zone:     portal.PortalStone.Zone,
	}
}

func randomWeapon(portal *sworld.Portal) sworld.Item {
	return &sworld.Weapon{
		Damage: 10,
	}
}

func randomItemEvent(portal *sworld.Portal) *sworld.PortalEvent {
	item := portal.PortalStone.Zone.DropItem(portal)

	return &sworld.PortalEvent{Item: item}
}

func createDefaultZone() *sworld.Zone {
	zone := sworld.NewZone("Forest")

	// TODO: I have no idea where to put this
	zone.AddItemDrop(0, 10, func(portal *sworld.Portal) sworld.Item {
		return portal.RandomPortalStone()
	})
	zone.AddItemDrop(2, 2, randomPowerStone)
	zone.AddItemDrop(3, 6, randomPowerStone)
	zone.AddItemDrop(1, 10, randomWeapon)

	zone.AddEventDrop(0, 10, func(portal *sworld.Portal, position int) *sworld.PortalEvent {
		return randomItemEvent(portal)
	})
	zone.AddEventDrop(0, 50, func(portal *sworld.Portal, position int) *sworld.PortalEvent {
		return nil
	})
	zone.AddEventDrop(1, 15, func(portal *sworld.Portal, position int) *sworld.PortalEvent {
		return randomItemEvent(portal)
	})
	zone.AddEventDrop(1, 30, func(portal *sworld.Portal, position int) *sworld.PortalEvent {
		return portal.RandomEnemyEvent(position)
	})
	zone.AddEventDrop(2, 40, func(portal *sworld.Portal, position int) *sworld.PortalEvent {
		return portal.RandomEnemyEvent(position)
	})
	return zone
}
