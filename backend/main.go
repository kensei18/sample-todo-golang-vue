package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienschmidt/httprouter"
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

func (t *Task) parse(r io.ReadCloser) {
	data, _ := ioutil.ReadAll(r)
	err := json.Unmarshal(data, t)
	if err != nil {
		panic(err)
	}

	if err = r.Close(); err != nil {
		panic(err)
	}
}

func (t *Task) find(db *pg.DB, id uint) {
	if err := db.Model(t).Where("id = ?", id).Select(); err != nil {
		log.Println(err)
	}
}

func (t *Task) create(db *pg.DB) {
	if _, err := db.Model(t).Insert(); err != nil {
		panic(err)
	}
}

func (t *Task) update(db *pg.DB) {
	if _, err := db.Model(t).WherePK().Update(); err != nil {
		log.Println(err)
	}
}

func (t *Task) delete(db *pg.DB) {
	if _, err := db.Model(t).WherePK().Delete(); err != nil {
		log.Println(err)
	}
}

type Tasks []Task

func (t *Tasks) get(db *pg.DB) {
	if err := db.Model(t).Select(); err != nil {
		panic(err)
	}
}

type databaseHandler func(*pg.DB)

func connectDatabase(f databaseHandler) {
	db := getDatabase()
	f(db)
	if err := db.Close(); err != nil {
		log.Println(err)
	}
}

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

func getTasksHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	tasks := new(Tasks)
	connectDatabase(func(db *pg.DB) {
		tasks.get(db)
	})

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}

	if _, err = w.Write(tasksJson); err != nil {
		panic(err)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var task Task
	connectDatabase(func(db *pg.DB) {
		task.parse(r.Body)
		task.create(db)
	})

	taskJson, _ := json.Marshal(task)
	if _, err := w.Write(taskJson); err != nil {
		panic(err)
	}
}

func updateTaskHandler(_ http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	task := new(Task)
	connectDatabase(func(db *pg.DB) {
		task.find(db, uint(id))
		task.parse(r.Body)
		task.update(db)
	})
}

func deleteTaskHandler(_ http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	id, _ := strconv.Atoi(p.ByName("id"))
	task := new(Task)
	connectDatabase(func(db *pg.DB) {
		task.find(db, uint(id))
		task.delete(db)
	})
}

func accessLog(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		h(w, r, p)
	}
}

func main() {
	router := httprouter.New()

	// /tasks
	router.GET("/api/tasks", accessLog(getTasksHandler))
	router.POST("/api/tasks", accessLog(createTaskHandler))
	router.PUT("/api/tasks/:id", accessLog(updateTaskHandler))
	router.DELETE("/api/tasks/:id", accessLog(deleteTaskHandler))

	log.Fatal(http.ListenAndServe(":5000", router))
}
