package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type ListItem struct {
	Id          int    `json:"id,string"`
	Checked     bool   `json:"checked,string"`
	Description string `json:"description"`
}

type List struct {
	Id          int        `json:"id,string"`
	Description string     `json:"description"`
	Items       []ListItem `json:"items"`
	Date        string     `json:"created_on"`
}

type Store struct {
	db *sql.DB
}

func (s *Store) listIndexHandler(w http.ResponseWriter, r *http.Request) {
	var lists []*List
	cutoff := time.Now().Add(-1 * time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	rows, err := s.db.Query("SELECT id, description, created_on FROM shopping_lists WHERE deleted=? AND created_on >= ?", 0, cutoff)
	if err != nil {
		fmt.Println("ERROR retrieving from DB - ", err)
	}

	for rows.Next() {
		list := &List{}
		rows.Scan(&list.Id, &list.Description, &list.Date)
		lists = append(lists, list)
	}
	rows.Close()

	sendResponseInJson(w, lists)
}

func (s *Store) listAddHandler(w http.ResponseWriter, r *http.Request) {
	list := &List{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&list)
	defer r.Body.Close()

	result, err := s.db.Exec("INSERT INTO shopping_lists(description) VALUES(?)", list.Description)
	if err != nil {
		fmt.Println("ERROR saving to DB - ", err)
	}
	id64, err := result.LastInsertId()
	newId := int(id64)

	list = &List{Id: newId}
	s.db.QueryRow("SELECT description, created_on FROM shopping_lists WHERE id = ?", list.Id).Scan(&list.Description, &list.Date)
	sendResponseInJson(w, list)
}

func (s *Store) listUpdateHandler(w http.ResponseWriter, r *http.Request) {

	listParams := &List{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&listParams)
	defer r.Body.Close()

	_, err = s.db.Exec("UPDATE shopping_lists SET description=? WHERE id=?", listParams.Description, listParams.Id)
	if err != nil {
		fmt.Println("ERROR saving to DB - ", err)
	}

	list := &List{Id: listParams.Id}
	err = s.db.QueryRow("SELECT description, created_on FROM shopping_lists WHERE id = ?", list.Id).Scan(&list.Description, &list.Date)
	sendResponseInJson(w, list)
}

func (s *Store) listGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("error error")
	}

	list := &List{Id: id}
	err = s.db.QueryRow("SELECT description, created_on FROM shopping_lists WHERE id = ? AND deleted=0", list.Id).Scan(&list.Description, &list.Date)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	sendResponseInJson(w, list)
}

func (s *Store) listDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s.db.Exec("UPDATE shopping_lists SET deleted=1 WHERE id=?", id)
	w.WriteHeader(200)
}

func (s *Store) listAddItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	data := List{Id: id}

	sendResponseInJson(w, data)
}

func (s *Store) listRemoveItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	data := List{Id: id}

	sendResponseInJson(w, data)
}

func (s *Store) listItemUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	data := List{Id: id}

	sendResponseInJson(w, data)
}

func sendResponseInJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, string(jsonData))
	}
}

func main() {

	db, err := sql.Open("mysql", "po:roach@tcp(127.0.0.1:3306)/posworkshop_dev")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := &Store{db: db}

	r := mux.NewRouter()

	r.HandleFunc("/lists", store.listIndexHandler).Methods("GET")
	r.HandleFunc("/lists", store.listAddHandler).Methods("POST")
	r.HandleFunc("/lists/{id}", store.listGetHandler).Methods("GET")
	r.HandleFunc("/lists/{id}", store.listUpdateHandler).Methods("PUT")
	r.HandleFunc("/lists/{id}", store.listDeleteHandler).Methods("DELETE")

	// r.HandleFunc("/lists/{id}/addItem", store.listAddItemHandler).Methods("PUT")
	// r.HandleFunc("/lists/{id}/RemoveItem", store.listRemoveItemHandler).Methods("PUT")
	// r.HandleFunc("/listitems/{id}", store.listItemUpdateHandler).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (li *ListItem) check()   {}
func (li *ListItem) uncheck() {}

func (l *List) add()                  {}
func (l *List) get()                  {}
func (l *List) addItem(li ListItem)   {}
func (l *List) removeItem(itemId int) {}
func (l *List) del()                  {}

// handlers
func listIndex() []List {
	return []List{}
}

func listGet(listId int) List {
	return List{}
}

func listAdd(l List) {}

func listAddItem(listId int, li ListItem)   {}
func listRemoveItem(listId int, itemId int) {}

func listDel(l List) {}

func listItemCheck(itemId int)   {}
func listItemUncheck(itemId int) {}
