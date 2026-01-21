package playerLogic

import "fmt"

type LevelingCurve struct{
	LevelCaps []int
	ExperienceRequired []int
}

var levelCaps = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

var fastLevelingCurve = []int{
	0, 1300, 3300, 6000, 10000,
	15000, 23000, 34000, 50000, 71000,
	100000, 140000, 200000, 290000, 410000, 550000,
	750000, 1000000, 1350000, 1800000,
}


var normalLevelingCurve = []int{
	0, 2000, 5000, 9000, 15000,
	23000, 35000, 50000, 75000, 105000,
	155000, 220000, 315000, 445000, 635000, 890000,
	1300000, 1800000, 2550000, 3600000,
}

var slowLevelingCurve = []int{
	0, 3000, 7500, 14000, 23000,
	35000, 53000, 77000, 115000, 160000,
	235000, 330000, 475000, 665000, 955000, 1350000,
	1900000, 2700000, 3850000, 5350000,
} 

type LevelingTrack string
const (
	LevelingTrackFast LevelingTrack = "fast"
	LevelingTrackNormal LevelingTrack = "normal"
	LevelingTrackSlow LevelingTrack = "slow"
)

var LevelingTracks = []LevelingTrack{
	LevelingTrackFast,
	LevelingTrackNormal,
	LevelingTrackSlow,
}

// GetLevelingCurve returns the leveling curve for players.
func GetLevelingCurve(track LevelingTrack) LevelingCurve {
	switch track {
		case LevelingTrackSlow:
			return LevelingCurve{
				LevelCaps: levelCaps,
				ExperienceRequired: slowLevelingCurve,
			}
		case LevelingTrackNormal:
			return LevelingCurve{
				LevelCaps: levelCaps,
				ExperienceRequired: normalLevelingCurve,
			}
		case LevelingTrackFast:
			return LevelingCurve{
				LevelCaps: levelCaps,
				ExperienceRequired: fastLevelingCurve,
			}
		default:
			return LevelingCurve{
				LevelCaps: levelCaps,
				ExperienceRequired: normalLevelingCurve,
			}
	}
}

// CanLevelUp checks if the player has enough experience to level up.
func CanLevelUp(player *Player) bool {
	curve := GetLevelingCurve(player.LevelTrack)
	if player.Level >= len(curve.ExperienceRequired) {
		return false
	}
	requiredXP := curve.ExperienceRequired[player.Level]
	return player.Experience >= requiredXP
}

// CanLevelUpTo checks if the player can level up to the target level.
func CanLevelUpTo(player *Player, targetLevel int) bool {
	curve := GetLevelingCurve(player.LevelTrack)
	if targetLevel <= player.Level || targetLevel > len(curve.ExperienceRequired) {
		return false
	}
	requiredXP := curve.ExperienceRequired[targetLevel-1]
	return player.Experience >= requiredXP
}

// CanLevelDown checks if the player can level down based on their experience when changing level tracks.
func CanLevelDown(player *Player) bool{
	curve := GetLevelingCurve(player.LevelTrack)
	if player.Level <= 1 || player.Level > len(curve.ExperienceRequired) {
		return false
	}
	requiredXP := curve.ExperienceRequired[player.Level-1]
	return player.Experience < requiredXP
}


// LevelUp increases the player's level and updates their stats accordingly.
func LevelUp(player *Player) {
	player.Level++
	player.SetAttributeModifiers()
	player.SetDerivedStats()
}

// LevelDown decreases the player's level and updates their stats accordingly.
func LevelDown(player *Player) {
	if player.Level > 1 {
		player.Level--
		player.SetAttributeModifiers()
		player.SetDerivedStats()
	}
}

// GetExperienceForNextLevel returns the experience required for the next level.
func GetExperienceForNextLevel(player *Player) int {
	curve := GetLevelingCurve(player.LevelTrack)
	if player.Level >= len(curve.ExperienceRequired) {
		return -1 // Max level reached
	}
	return curve.ExperienceRequired[player.Level]
}

// GetExperienceForLevel returns the experience required for a specific level.
func GetExperienceForLevel(level int, track LevelingTrack) int {
	curve := GetLevelingCurve(track)
	if level <= 0 || level > len(curve.ExperienceRequired) {
		return -1 // Invalid level
	}
	return curve.ExperienceRequired[level-1]
}

// GetLevelTrack converts a string to a LevelingTrack type.
func GetLevelTrack(track string) (LevelingTrack, error){
	switch track {
		case "fast":
			return LevelingTrackFast, nil
		case "normal":
			return LevelingTrackNormal, nil
		case "slow":
			return LevelingTrackSlow, nil
		default:
			return LevelingTrackNormal, fmt.Errorf("invalid leveling track: %s", track)
	}
}
