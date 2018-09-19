package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	mux "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var idCounter int

type MyHandler struct {
	mutex sync.Mutex
}

type Shoe struct {
	ID    string
	Model string
	Brand string
	Price float32
}

func handlerWrapper(h http.HandlerFunc) http.HandlerFunc {
	return basicAuth(h)
}

func (m *MyHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 HTTP status code returned!"))
}

func generateID(shoe Shoe) string {
	idCounter++
	var str []string
	str = append(str, strings.ToLower(shoe.Brand[:3]))
	str = append(str, strconv.Itoa(idCounter))
	return strings.Join(str, "-")
}

func (m *MyHandler) Create(w http.ResponseWriter, r *http.Request) {
	shoe := Shoe{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&shoe); err != nil {
		log.Info(err)
		fmt.Errorf("Error json decode %v", err)
		return
	}
	defer r.Body.Close()
	shoe.ID = generateID(shoe)

	if err := db.Save(&shoe).Error; err != nil {
		log.Info("Error save")
		fmt.Errorf("Error save ind db %v", err)
		return
	}
	log.Info("Added shoe")
	fmt.Fprintln(w, "Succesfully added!")
}

func (m *MyHandler) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	shoe := Shoe{ID: id}
	db.Delete(&shoe)
	fmt.Fprintln(w, "Succesfully deleted!")
}

func (m *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	var shoe Shoe
	vars := mux.Vars(r)
	db.Where("id = ?", vars["id"]).First(&shoe)
	log.Info("here it is:", vars["id"], shoe)
}

func (m *MyHandler) List(w http.ResponseWriter, r *http.Request) {
	var list []Shoe
	db.Find(&list)

	listJson, err := json.Marshal(list)
	if err != nil {
		log.Info(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(listJson)
}

func main() {
	var err error
	dbUrl := os.Getenv("DB_URL")
	dbDialect := os.Getenv("DB_DIALECT")
	db, err = gorm.Open(dbDialect, dbUrl)
	if err != nil {
		panic("Failed to connect DB")
	}
	defer db.Close()
	db.AutoMigrate(&Shoe{})

	h := MyHandler{}

	r := mux.NewRouter()
	r.HandleFunc("/", handlerWrapper(h.StatusHandler))
	r.HandleFunc("/shoe", handlerWrapper(h.Create)).Methods("POST")
	r.HandleFunc("/shoe/{id:[a-z]+-[0-9]+}", handlerWrapper(h.Remove)).Methods("DELETE")
	r.HandleFunc("/shoe/{id:[a-z]+-[0-9]+}", handlerWrapper(h.Get)).Methods("GET")
	r.HandleFunc("/shoe", handlerWrapper(h.List)).Methods("GET")
	port := os.Getenv("SHOE_PORT")
	log.Fatal(http.ListenAndServe(port, r))
}
