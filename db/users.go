package db

type User struct {
	Uid       int    //`json:"uid"`
	Username  string //`json:"username"`
	Password  string //`json:"password"`
	PilotName string //`json:"pilotName"`
	Level     int    //`json:"level"`
	Rank      string //`json:"rank"`
	Credits   int    //`json:"credits"`
}
