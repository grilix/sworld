package server

import "github.com/grilix/sworld/sworld"

// UserDetails represents a user
type UserDetails struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// CharacterDetails represents a character in a response
type CharacterDetails struct {
	ID        string `json:"id"`
	Level     int    `json:"level"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"max_health"`
	Exploring bool   `json:"exploring"`
}

// StoneDetails holds the information about a stone item in a response
type StoneDetails struct {
	Level    int         `json:"level"`
	Zone     ZoneDetails `json:"zone"`
	Duration string      `json:"duration"`
}

// BagSlotDetails represents a bag slot on a response
type BagSlotDetails struct {
	// TODO:
	Slot  int           `json:"slot"`
	Item  string        `json:"item,omitempty"`
	Stone *StoneDetails `json:"stone,omitempty"`
}

// BagDetails represents a bag in a response
type BagDetails struct {
	ID    int               `json:"id"`
	Items []*BagSlotDetails `json:"items"`
}

// ZoneDetails represents the details of a zone
type ZoneDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PortalDetails holds the details for a portal
type PortalDetails struct {
	ID       string         `json:"id"`
	IsOpen   bool           `json:"is_open"`
	Duration int            `json:"duration"`
	TimeLeft int            `json:"time_left"`
	Level    int            `json:"level"`
	Zone     ZoneDetails    `json:"zone"`
	Enemies  []EnemyDetails `json:"enemies,omitempty"`
}

// EnemyDetails contains details about an enemy
type EnemyDetails struct {
	ID        string `json:"id"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"max_health"`
	Level     int    `json:"level"`
}

// TODO: Not sure what could be the best approach here, we'll just try every item type for now
func bagSlotDetails(slot int, item sworld.Item) *BagSlotDetails {
	details := &BagSlotDetails{
		Slot: slot,
	}

	stoneItem, ok := item.(*sworld.PortalStone)
	if ok {
		details.Item = "stone"
		details.Stone = &StoneDetails{
			Level:    stoneItem.Level,
			Duration: stoneItem.Duration.String(),
			Zone: ZoneDetails{
				ID:   stoneItem.Zone.ID,
				Name: stoneItem.Zone.Name,
			},
		}
		return details
	}
	return details
}

func inventoryDetails(inventory []sworld.Bag) []*BagDetails {
	characterBags := make([]*BagDetails, 0, len(inventory))
	for id, bag := range inventory {
		items := bag.Items()

		bagItems := make([]*BagSlotDetails, 0, len(items))
		for slot, item := range items {
			bagItems = append(bagItems, bagSlotDetails(slot, item))
		}

		characterBags = append(characterBags, &BagDetails{
			ID:    id,
			Items: bagItems,
		})
	}
	return characterBags
}

func portalDetails(portal *sworld.Portal, listing bool) *PortalDetails {
	var deadEnemies []EnemyDetails

	if !listing {
		enemies := portal.DeadEnemies()
		deadEnemies = make([]EnemyDetails, 0, len(enemies))
		for _, enemy := range enemies {
			deadEnemies = append(deadEnemies, EnemyDetails{
				ID:        enemy.ID,
				Health:    enemy.Health,
				MaxHealth: enemy.MaxHealth,
				Level:     enemy.Level,
			})
		}
	}

	timeLeft := 0
	if portal.IsOpen {
		timeLeft = int(portal.TimeLeft().Seconds())
	}
	return &PortalDetails{
		ID:       portal.ID,
		IsOpen:   portal.IsOpen,
		Duration: int(portal.PortalStone.Duration.Seconds()),
		TimeLeft: timeLeft,
		Level:    portal.PortalStone.Level,
		Enemies:  deadEnemies,
		Zone: ZoneDetails{
			ID:   portal.PortalStone.Zone.ID,
			Name: portal.PortalStone.Zone.Name,
		},
	}
}
