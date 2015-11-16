/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/andysctu/iMND2/services"
	_ "github.com/lib/pq"
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

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "user=andy dbname=testDB sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	db := initDB()
	services.InitDBSvc(db)

	http.HandleFunc("/api/comments", handleComments)
	http.HandleFunc("/contactInfo", handleContactInfo)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
