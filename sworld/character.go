package sworld

// Killing enemy xp:
// xp := int64(
//     math.Round(
//         (float64(enemy.Level) / float64(c.Level)) * float64(enemy.Level),
//     ),
// )

// Next lvl xp:
// return int64(math.Round((4 * math.Pow(float64(c.Level), 3)) / 5))

type Character struct {
	ID        string
	Health    int
	MaxHealth int
	Gold      int
	Exploring bool
	Level     int

	experience int64
}
