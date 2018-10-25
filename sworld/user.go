package sworld

type User struct {
	ID        string
	Username  string
	Character *Character

	// TODO: inventory
	Gold,
	Items,
	Stones,
	Materials int
}
