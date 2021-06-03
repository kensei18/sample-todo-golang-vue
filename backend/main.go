package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Task struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

func (t Task) String() string {
	return fmt.Sprintf(
		"<Id: %d, Name: %s, Description: %s, Status: %d>",
		t.Id, t.Name, t.Description, t.Status,
	)
}

func (t *Task) parseBody(r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, t)
	if err != nil {
		panic(err)
	}
}

func (t *Task) find(db *pg.DB, id uint) {
	if err := db.Model(t).Where("id = ?", id).Select(); err != nil {
		log.Println(err)
	}
}

var validPath = regexp.MustCompile("^/api/(tasks)/([0-9]*)$")

func getDatabase() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		Password: "password",
		Database: "sample_todo_golang_vue",
	})

	_ = createSchema(db)
	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Task)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	db := getDatabase()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		getTasks(w, db)
	case "POST":
		createTask(w, r, db)
	case "PUT":
		updateTask(r, m, db)
	case "DELETE":
		deleteTask(m, db)
	}
}

func getTasks(w http.ResponseWriter, db *pg.DB) {
	var tasks []Task
	if err := db.Model(&tasks).Select(); err != nil {
		panic(err)
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}

	if _, err = w.Write(tasksJson); err != nil {
		panic(err)
	}
}

func createTask(w http.ResponseWriter, r *http.Request, db *pg.DB) {
	var task Task
	task.parseBody(r)
	if _, err := db.Model(&task).Insert(); err != nil {
		panic(err)
	}

	taskJson, _ := json.Marshal(task)
	if _, err := w.Write(taskJson); err != nil {
		panic(err)
	}
}

func updateTask(r *http.Request, m []string, db *pg.DB) {
	id, _ := strconv.Atoi(m[2])
	task := new(Task)
	task.find(db, uint(id))

	task.parseBody(r)
	if _, err := db.Model(task).WherePK().Update(); err != nil {
		log.Println(err)
	}
}

func deleteTask(m []string, db *pg.DB) {
	id, _ := strconv.Atoi(m[2])
	task := new(Task)
	task.find(db, uint(id))
	if _, err := db.Model(task).WherePK().Delete(); err != nil {
		log.Println(err)
	}
}

func accessLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		h(w, r)
	}
}

func main() {
	http.HandleFunc("/api/tasks/", accessLog(tasksHandler))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
