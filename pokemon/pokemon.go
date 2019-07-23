// Package pokemon provides objects and methods representing a Pokemon in Pokemon Go. It does not
// provide a "New" function: to load from a gamemaster.json, use the gamemaster package's PokemonBy.* methods.
package pokemon

import (
	"math"
)

// CPM is the combat multiplier used when calculating CP. Each entry is a half level: index 0 is
// level 1, index 1 is 1.5, 2 is 2, etc.
var CPM = []float64{
	0.094,
	0.1351374318,
	0.16639787,
	0.192650919,
	0.21573247,
	0.2365726613,
	0.25572005,
	0.2735303812,
	0.29024988,
	0.3060573775,
	0.3210876,
	0.3354450362,
	0.34921268,
	0.3624577511,
	0.3752356,
	0.387592416,
	0.39956728,
	0.4111935514,
	0.4225,
	0.4329264091,
	0.44310755,
	0.4530599591,
	0.4627984,
	0.472336093,
	0.48168495,
	0.4908558003,
	0.49985844,
	0.508701765,
	0.51739395,
	0.5259425113,
	0.5343543,
	0.5426357375,
	0.5507927,
	0.5588305862,
	0.5667545,
	0.5745691333,
	0.5822789,
	0.5898879072,
	0.5974,
	0.6048236651,
	0.6121573,
	0.6194041216,
	0.6265671,
	0.6336491432,
	0.64065295,
	0.6475809666,
	0.65443563,
	0.6612192524,
	0.667934,
	0.6745818959,
	0.6811649,
	0.6876849038,
	0.69414365,
	0.70054287,
	0.7068842,
	0.7131691091,
	0.7193991,
	0.7255756136,
	0.7317,
	0.7347410093,
	0.7377695,
	0.7407855938,
	0.74378943,
	0.7467812109,
	0.74976104,
	0.7527290867,
	0.7556855,
	0.7586303683,
	0.76156384,
	0.7644860647,
	0.76739717,
	0.7702972656,
	0.7731865,
	0.7760649616,
	0.77893275,
	0.7817900548,
	0.784637,
	0.7874736075,
	0.7903,
}

// Pokemon represents a Pokemon, including its stats. It can be created from a gamemaster.json file
// or saved -- via Save() -- into yaml.
type Pokemon struct {
	ID              string `json:"speciesId"`
	Dex             int    `json:"dex"`
	Name            string `json:"speciesName"`
	BaseStats       Stats  `json:"baseStats"`
	IVs             Stats
	CalculatedStats Stats
	Level           float64
	CP              int
	Moves           Moves
}

// Stats represent the three statistics of a Pokemon: attack, defense, and stamina/HP. This is used
// both for a Pokemon's base stats -- which every Pokemon of a type shares -- and its individual
// values, or IVs, which are unique.
type Stats struct {
	Attack  float64 `json:"atk"`
	Defense float64 `json:"def"`
	HP      float64 `json:"hp"`
}

// Moves are a Pokemon's moves. They can have one fast and one or two charge.
type Moves struct {
	Fast   string    `yaml:"Fast"`
	Charge [2]string `yaml:"Charge"`
}

// Calculate will use the level and IVs of a Pokemon to calculate its CP and final stats. If they
// have already been calculated this will be a no-op.
//
// CalculatedStats are saved separately, as it's useful to know the stat product of a Pokemon. The
// calculated stats include the rounded down HP and as such is different from what's used to calculate
// CP.
func (p *Pokemon) Calculate() {
	if p.CP > 0 {
		return
	}

	level := int((p.Level - 1) * 2)
	attack := (p.BaseStats.Attack + p.IVs.Attack) * CPM[level]
	defense := (p.BaseStats.Defense + p.IVs.Defense) * CPM[level]
	stamina := (p.BaseStats.HP + p.IVs.HP) * CPM[level]

	cp := math.Floor((math.Pow(stamina, 0.5) * attack * math.Pow(defense, 0.5)) / 10)
	if cp < 10 {
		cp = 10
	}
	p.CP = int(cp)

	p.CalculatedStats = Stats{
		Attack:  attack,
		Defense: defense,
		HP:      math.Floor(stamina),
	}
}
