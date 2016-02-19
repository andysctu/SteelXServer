package db

type Mech struct {
	Uid       int    `json:"uid"`
	Arms      string `json:"arms"`
	Legs      string `json:"legs"`
	Core      string `json:"core"`
	Head      string `json:"head"`
	Weapon1L  int    `json:"weapon1l"`
	Weapon1R  int    `json:"weapon1r"`
	Weapon2L  int    `json:"weapon2l"`
	Weapon2R  int    `json:"weapon2r"`
	Booster   int    `json:"booster"`
	IsPrimary bool   `json:"isPrimary"`
}
