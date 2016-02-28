package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/andysctu/iMND2/Godeps/_workspace/src/github.com/lib/pq"
	mydb "github.com/andysctu/iMND2/db"
	"github.com/andysctu/iMND2/services"
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

type comment struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

const dataFile = "./comments.json"
const contactInfoFile = "./contactInfo.json"

var commentMutex = new(sync.Mutex)

// Handle comments
func handleComments(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	commentMutex.Lock()
	defer commentMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	// Read the comments from the file.
	commentData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var comments []comment
		if err := json.Unmarshal(commentData, &comments); err != nil {
			http.Error(w, fmt.Sprintf("Unable to Unmarshal comments from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		// Add a new comment to the in memory slice of comments
		comments = append(comments, comment{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})

		// Marshal the comments to indented json.
		commentData, err = json.MarshalIndent(comments, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal comments to json: %s", err), http.StatusInternalServerError)
			return
		}

		// Write out the comments to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, commentData, fi.Mode())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to write comments to data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		io.Copy(w, bytes.NewReader(commentData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(commentData))

	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

// Handle comments
func handleContactInfo(w http.ResponseWriter, r *http.Request) {
	// Stat the file, so we can find its current permissions
	// fi, err := os.Stat(contactInfoFile)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
	// 	return
	// }

	// Read the comments from the file.
	contactInfo, err := ioutil.ReadFile(contactInfoFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", contactInfoFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":

	case "GET":
		var contactInfoMap map[string]string
		err = json.Unmarshal(contactInfo, &contactInfoMap)
		if err != nil {
			log.Println(err)
		}
		SendResponse(w, 200, contactInfoMap)
	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	success := false
	ret := make(map[string]interface{})
	db := services.GetDB()
	var user mydb.User
	switch r.Method {
	case "POST":
		{
			log.Println("1")
			potentialPassword := r.FormValue("password")

			rows, err := db.Query("SELECT * FROM users where username = ?", r.FormValue("username")) // where ... sql injection
			if err != nil {
				log.Fatal(err)
			}

			log.Println("2")
			if rows.Next() {
				err = rows.Scan(
					&user.Uid,
					&user.Username,
					&user.Password,
					&user.PilotName,
					&user.Level,
					&user.Rank,
					&user.Credits,
				)
				if err != nil {
					// log.Println("here2")
					log.Fatal(err)
				}
				log.Println("3")
				log.Println("potpass: " + potentialPassword)
				log.Println("userpas: " + user.Password)
				if potentialPassword == user.Password {
					log.Println("4")
					w.WriteHeader(http.StatusOK)
					log.Println("success")
					success = true
					ret["User"] = user
				}
			} else {
				log.Println("5")
				log.Println("failure")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			rows.Close()
			log.Println("6")
			if success {
				log.Println("7")
				rows, err := db.Query("SELECT * FROM mechs WHERE uid = $1 AND isPrimary = true", user.Uid) // sql injection

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
					log.Println("8")
					log.Println(mech)

					ret["Mech"] = mech

					SendResponse(w, http.StatusOK, ret)
				} else {
					log.Println("9")
					w.WriteHeader(http.StatusNotFound)
					log.Println("No mech data for user: " + string(user.Uid))
				}
				log.Println("10")
			}
		}
	}
}

func MechHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			uid := r.FormValue("uid")
			log.Println(uid)
			db := services.GetDB()

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

				log.Println(mech)

				SendResponse(w, http.StatusOK, mech)
				return

			}

			w.WriteHeader(http.StatusNotFound)

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
		port = "3000"
	}

	db := initDB()
	services.InitDBSvc(db)

	http.HandleFunc("/api/comments", handleComments)
	http.HandleFunc("/contactInfo", handleContactInfo)
	http.HandleFunc("/mech", MechHandler)
	http.HandleFunc("/login", LoginHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
