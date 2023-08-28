package constants

import "strings"

type Mode int

const (
	TeamDeathmatch Mode = iota
	BattleRoyale
	GunSmith
	OneVsOne
	Mayhem
	Unknown
)

func ParseMode(modeStr string) Mode {
	switch strings.ToLower(modeStr) {
	case "team deathmatch":
		return TeamDeathmatch
	case "battle royale":
		return BattleRoyale
	case "gunsmith":
		return GunSmith
	case "1 v 1":
		return OneVsOne
	case "mayhem":
		return Mayhem
	default:
		return Unknown
	}
}

// Defined a map to store the room limit for each mode
var roomLimitMap = map[Mode]int{
	TeamDeathmatch: 10,
	BattleRoyale:   20,
	GunSmith:       8,
	OneVsOne:       2,
	Mayhem:         5,
}

// Function to get the max players for a given mode
func RoomLimit(mode Mode) int {
	return roomLimitMap[mode]
}
