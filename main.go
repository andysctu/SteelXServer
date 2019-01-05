package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/andysctu/SteelXServer/Godeps/_workspace/src/github.com/lib/pq"
	mydb "github.com/andysctu/SteelXServer/db"
	// mydb "./db"
	"github.com/andysctu/SteelXServer/services"
	// "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

func SendResponse(w http.ResponseWriter, status int, resp interface{}) {
	val := reflect.ValueOf(resp)

	if err, ok := val.Interface().(error); ok {
		SendStringResponse(w, status, err.Error())
		return
	}

	fmt.Printf("Sending: %v\n", resp)
	switch resp := resp.(type) {
	case string:
		SendStringResponse(w, status, resp)
	default:
		bytes, _ := json.Marshal(resp)
		SendJsonResponse(w, status, string(bytes))
	}

}

func SendJsonResponse(w http.ResponseWriter, status int, resp string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, resp)

}

func SendStringResponse(w http.ResponseWriter, status int, str string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, str)
}

// func PilotInfoHandler(w http.ResponseWriter, r *http.Request) {
// 	// ret := make(map[string]interface{})
// 	// db := services.GetDB()
// 	uid := r.FormValue("uid")
// 	log.Println(uid)
// 	switch r.Method {
// 	case "PUT":
// 		{
// 			for k, v := range r.Form {
// 				// db.Exec("")
// 				log.Println(k)
// 				log.Println(v)
// 			}
// 		}
// 	}
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	success := false
	ret := make(map[string]interface{})
	db := services.GetDB()
	var user mydb.User
	switch r.Method {
	case "POST":
		{
			potentialPassword := r.FormValue("password")

			rows, err := db.Query("SELECT * FROM users WHERE username = $1", r.FormValue("username")) // where ... sql injection
			if err != nil {
				log.Fatal(err)
			}
			log.Println(potentialPassword)
			if rows.Next() {
				err = rows.Scan(
					&user.Uid,
					&user.Username,
					&user.Password,
					&user.PilotName,
					&user.Level,
					&user.Rank,
					&user.Credits,
					&user.Kills,
					&user.Deaths,
					&user.Assists,
					&user.TimeLogged,
				)
				if err != nil {
					log.Fatal(err)
				}

				if potentialPassword == user.Password {
					success = true
					ret["User"] = user
				}
			}
			rows.Close()

			if !success {
				log.Printf("Invalid credentials")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Get main Mech
			rows, err = db.Query("SELECT * FROM mechs WHERE uid = $1 AND isPrimary = true", user.Uid) // sql injection

			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}

			if rows.Next() {
				var mech mydb.Mech
				err = rows.Scan(
					&mech.Uid,
					&mech.Arms,
					&mech.Legs,
					&mech.Core,
					&mech.Head,
					&mech.Weapon1L,
					&mech.Weapon1R,
					&mech.Weapon2L,
					&mech.Weapon2R,
					&mech.Booster,
					&mech.IsPrimary,
				)
				if err != nil {
					log.Fatal(err)
				}

				ret["Mech"] = []mydb.Mech{mech}

			} else {
				w.WriteHeader(http.StatusNotFound)
				log.Println("No mech data for user: " + string(user.Uid))
				return
			}

			// Get all equipment owned
			ret["Owns"] = make([]string, 0)
			rows, err = db.Query("SELECT name FROM equipment E, owns O, users U WHERE E.eid = O.eid and O.uid = U.uid and U.uid = $1;", user.Uid)
			for rows.Next() {
				var part string
				err = rows.Scan(&part)
				if err != nil {
					log.Fatal(err)
				}
				ret["Owns"] = append(ret["Owns"].([]string), part)
			}

			SendResponse(w, http.StatusOK, ret)

		}
	}
}

func MechHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	log.Println(uid)
	db := services.GetDB()

	switch r.Method {
	case "GET":
		{

			rows, err := db.Query("SELECT * FROM mechs WHERE uid = $1 AND isPrimary = true", uid)

			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				var mech mydb.Mech
				err = rows.Scan(
					&mech.Uid,
					&mech.Arms,
					&mech.Legs,
					&mech.Core,
					&mech.Head,
					&mech.Weapon1L,
					&mech.Weapon1R,
					&mech.Weapon2L,
					&mech.Weapon2R,
					&mech.Booster,
					&mech.IsPrimary,
				)
				if err != nil {
					log.Fatal(err)
				}
				SendResponse(w, http.StatusOK, mech)
				return

			}

			w.WriteHeader(http.StatusNotFound)

		}
	case "POST":
		{
			for k, v := range r.PostForm {
				log.Println(k)
				log.Println(v)

				if k == "uid" {
					continue
				}

				// Need to check if they own it
				_, err := db.Exec(fmt.Sprintf("UPDATE mechs SET %s = $1 WHERE uid = $2", k), v[0], uid)
				if err != nil {
					log.Println(err)
				}
			}

			rows, err := db.Query("SELECT * FROM mechs WHERE uid = $1 AND isPrimary = true", uid)

			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				var mech mydb.Mech
				err = rows.Scan(
					&mech.Uid,
					&mech.Arms,
					&mech.Legs,
					&mech.Core,
					&mech.Head,
					&mech.Weapon1L,
					&mech.Weapon1R,
					&mech.Weapon2L,
					&mech.Weapon2R,
					&mech.Booster,
					&mech.IsPrimary,
				)
				if err != nil {
					log.Fatal(err)
				}
				SendResponse(w, http.StatusOK, mech)
				return

			}

			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	db := services.GetDB()

	switch r.Method {
	case "POST":
		{
			// Get credit balance
			var credits int

			rows, err := db.Query("SELECT credits FROM users WHERE uid = $1", uid)
			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}

			if rows.Next() {
				err = rows.Scan(&credits)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Get equipment info
			eid := r.FormValue("eid")
			rows, err = db.Query("SELECT * FROM equipment WHERE eid = $1", eid)

			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}

			var equipment mydb.Equipment
			if rows.Next() {
				err = rows.Scan(
					&equipment.Eid,
					&equipment.Cost,
					&equipment.Type,
					&equipment.Name,
				)
				if err != nil {
					log.Fatal(err)
				}
			}

			log.Printf("credits: %d, cost: %d\n", credits, equipment.Cost)

			// Check if balance is sufficient
			if credits < equipment.Cost {
				SendResponse(w, http.StatusOK, false)
				return
			}

			// Begin atomic db write
			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				SendResponse(w, http.StatusOK, false)
				return
			}

			// Deduct credits
			credits -= equipment.Cost

			// Add user ownership of purchased item
			_, err = tx.Exec("INSERT INTO owns VALUES($1, $2)", uid, eid)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				SendResponse(w, http.StatusOK, false)
				return
			}

			_, err = tx.Exec("UPDATE users SET credits = $1 WHERE uid = $2", credits, uid)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				SendResponse(w, http.StatusOK, false)
				return
			}

			tx.Commit()

			SendResponse(w, http.StatusOK, true)
		}
	}
}

func GameHistoryHandler(w http.ResponseWriter, r *http.Request) {
	db := services.GetDB()
	// Need start time, end time, game type, victor
	// time format: MM/DD/YYYY HH:mm:SS
	// Need uid, kda of each player, team of each player
	switch r.Method {
	case "POST":
		{
			start_time := r.FormValue("start_time")
			end_time := r.FormValue("end_time")
			game_type := r.FormValue("game_type")
			victor := r.FormValue("victor")

			layout := "01/02/2006 15:04:05"
			t1, err := time.Parse(layout, start_time)
			t2, err := time.Parse(layout, end_time)

			if err != nil {
				log.Println(err)
			}

			durationInSeconds := t2.Sub(t1).Seconds()

			// Record game in game_history table
			_, err = db.Exec("INSERT INTO game_history (start_time, end_time, game_type, victor) VALUES($1, $2, $3, $4)", start_time, end_time, game_type, victor)
			if err != nil {
				log.Println(err)
				SendResponse(w, http.StatusOK, false)
				return
			}

			// Doesn't work for some reason
			// gid, err := result.LastInsertId()
			// if err != nil {
			// 	log.Println(err)
			// }

			// Get gid
			rows, err := db.Query("SELECT max(gid) FROM game_history")
			if err != nil {
				log.Println(err)
			}

			var gid int
			if rows.Next() {
				err = rows.Scan(&gid)
				if err != nil {
					log.Println(err)
				}
			}

			// Record individual records in player_history table
			raw_player_histories := r.FormValue("player_histories")

			byt := []byte(raw_player_histories)
			var player_histories map[int]interface{}
			if err := json.Unmarshal(byt, &player_histories); err != nil {
				panic(err)
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
			}
			for uid, history := range player_histories {
				history_map := history.(map[string]interface{})
				kills := history_map["kills"].(float64)
				deaths := history_map["deaths"].(float64)
				assists := history_map["assists"].(float64)
				_, err = tx.Exec("INSERT INTO player_history (gid, uid, kills, deaths, assists, team) VALUES($1, $2, $3, $4, $5, $6)",
					gid, int64(uid), kills, deaths, assists, history_map["team"].(string))
				if err != nil {
					log.Println(err)
				}

				// Update users personal k/d/a and time_logged
				fmt.Printf(
					"UPDATE users SET kills = kills + %f, deaths = deaths + %f, assists = assists + %f, time_logged = time_logged + %f WHERE uid = $1",
					kills, deaths, assists, durationInSeconds)
				_, err = tx.Exec(
					fmt.Sprintf(
						"UPDATE users SET kills = kills + %f, deaths = deaths + %f, assists = assists + %f, time_logged = time_logged + %f WHERE uid = $1",
						kills, deaths, assists, durationInSeconds), uid)
				if err != nil {
					log.Println(err)
				}
			}

			tx.Commit()

			SendResponse(w, http.StatusOK, true)
			return
		}
	}
}

func initDB() *sql.DB {
	url := "postgres://syuanntjlkwjoo:bPkYjz9Q4EUj4_U3rSniAH7ILr@ec2-54-83-53-120.compute-1.amazonaws.com:5432/djk4n55d220oe"
	// url := os.Getenv("DATABASE_URL") + "?sslmode=require"
	log.Println("DB_URL: " + url)
	db, err := sql.Open("postgres", url)
	// db, err := sql.Open("postgres", testURL)
	if err != nil {
		log.Fatal("Error connecting to db: " + err.Error())
	}
	return db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	db := initDB()
	services.InitDBSvc(db)

	http.HandleFunc("/mech", MechHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/purchase", PurchaseHandler)
	http.HandleFunc("/game_history", GameHistoryHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
