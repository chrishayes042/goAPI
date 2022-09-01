package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Star struct {
	ID			uint
	Name		string`gorm:"unique"`
	Description	string
	URL 		string
}
type App struct{
	DB *gorm.DB
}

func (a *App) Initialize(dbDrive string, dbURI string){
	db, err := gorm.Open(dbDrive, dbURI)
	if err != nil{
		panic("failed to connect database")
	}
	a.DB = db

	a.DB.AutoMigrate(&Star{})
}

func (a *App) ListHandler(w http.ResponseWriter, r *http.Request){
	var stars []Star

	a.DB.Find(&stars)
	starsJSON, _ := json.Marshal(stars)

	w.WriteHeader(200)
	w.Write([]byte(starsJSON))
}

func (a *App) handler(w http.ResponseWriter, r *http.Request){

	a.DB.Create(&Star{Name: "test"})

	var star Star
	a.DB.First(&star, "name = ?", "test")


	w.WriteHeader(200)
	w.Write([]byte("Hello world!"))

	a.DB.Delete(&star)
}

func main(){

	a := &App{}
	a.Initialize("sqlite3", "test.db")

	r := mux.NewRouter()

	r.HandleFunc("/stars", a.ListHandler).Methods("GET")

	http.HandleFunc("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	defer a.DB.Close()
}