package db

type Mech struct {
	Uid       int    //`json:"uid"`
	Arms      string //`json:"arms"`
	Legs      string //`json:"legs"`
	Core      string //`json:"core"`
	Head      string //`json:"head"`
	Weapon1L  string //`json:"weapon1l"`
	Weapon1R  string //`json:"weapon1r"`
	Weapon2L  string //`json:"weapon2l"`
	Weapon2R  string //`json:"weapon2r"`
	Booster   string //`json:"booster"`
	IsPrimary bool   //`json:"isPrimary"`
}
